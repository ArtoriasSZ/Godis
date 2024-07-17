package parser

import (
	"Godis/interface/resp"
	"Godis/lib/logger"
	"Godis/resp/reply"
	"bufio"
	"errors"
	"io"
	"runtime/debug"
	"strconv"
	"strings"
)

// PayLoad chan有效载荷
type PayLoad struct {
	Date resp.Reply // 客户端发送的和服务端发送的都是reply
	Err  error
}

// readState 解析器
type readState struct {
	readingMultiLine  bool     // 正在解析多行数据还是单行
	msgType           byte     // 用户消息类型
	expectedArgsCount int      // 期待解析的参数个数
	bulkLen           int64    // 指令长度
	args              [][]byte // 用户传参
}

// isFinished 判断是否解析完成
func (r readState) isFinished() bool {
	// 读的参数大于0，并且读的参数的长度==期望读到的参数数量那么完成了
	return len(r.args) > 0 && len(r.args) == r.expectedArgsCount
}

// ParseStream 使用协程解析字节流，外部通过管道读取消息
// 这个是对于一个用户的
func ParseStream(reader io.Reader) <-chan *PayLoad {
	ch := make(chan *PayLoad)
	go parse0(reader, ch)
	return ch
}

// *2\r\n$3\r\nser\r\n$3\r\nser\r\n
func parse0(reader io.Reader, ch chan *PayLoad) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	bufReader := bufio.NewReader(reader)
	state := readState{}
	var err error
	var msg []byte
	for true {
		var ioErr bool
		msg, ioErr, err = readline(bufReader, &state)
		// 这里没有读到东西就直接报错了，这不对吧
		if err != nil {
			// 连接关闭
			// io错误关闭管道，终止服务
			if ioErr {
				ch <- &PayLoad{
					Err: err,
				}
				close(ch)
				return
			}
			// 协议错误
			ch <- &PayLoad{
				Err: err,
			}
			state = readState{}
			continue
		}
		// 多行解析模式
		if !state.readingMultiLine {
			// 不是多行处理模式，先看看是不是没有初始化
			if msg[0] == '*' {
				err = parseMultiBulkHeader(msg, &state)
				// 解析失败，协议错误往上抛就是了
				if err != nil {
					ch <- &PayLoad{
						Err: err,
					}
					state = readState{}
					continue
				}
				// 期待传参是0，说明是空，本次解析结束下一次
				if state.expectedArgsCount == 0 {
					ch <- &PayLoad{
						Date: reply.NewEmptyMultiBulkReply(),
					}
					state = readState{}
					continue
				}
			} else if msg[0] == '$' {
				// 如果是一个字符串
				err = parseBulkHeader(msg, &state)
				if err != nil {
					ch <- &PayLoad{
						Err: err,
					}
					state = readState{}
					continue
				}
				if state.bulkLen == -1 {
					ch <- &PayLoad{
						Date: reply.NewNullBulkReply(),
					}
					state = readState{}
					continue
				}
			} else { // 正常，错误，数字的情况
				var lineReply resp.Reply
				lineReply, err = parseSingleLineReply(msg)
				ch <- &PayLoad{
					Date: lineReply,
					Err:  err,
				}
				state = readState{}
				continue
			}
		} else {
			err = readBody(msg, &state)
			if err != nil {
				ch <- &PayLoad{
					Err: err,
				}
				state = readState{}
				continue
			}
			if state.isFinished() {
				var result resp.Reply
				if state.msgType == '*' {
					result = reply.NewMultiBulkReply(state.args)
				} else if state.msgType == '$' {
					result = reply.NewBulkReply(state.args[0])
				}
				ch <- &PayLoad{
					Date: result,
					Err:  err,
				}
				state = readState{}
			}
		}
	}
}

// 只负责根据\r\n读取一行数据传递给msg，对于msg的处理逻辑在后面
func readline(bufRead *bufio.Reader, state *readState) ([]byte, bool, error) {
	var msg []byte
	var err error
	if state.bulkLen == 0 {
		// 初始化 即 读前面的数字
		// 可能读到中间有\n，这个时候并不应该返回错误，
		// 但是数字中间也不应该有\n,发送了直接error也没问题
		// 阻塞读，并不会返回io.EOF
		msg, err = bufRead.ReadBytes('\n')
		if err != nil {
			return nil, true, err
			// 是io错误
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol error:" + string(msg))
		}
	} else {
		// 读字符串
		// 这个时候中间读到了\r
		msg = make([]byte, state.bulkLen+2)
		// ReadFull把io塞满切片中，下次再读重上次的位置继续
		_, err = io.ReadFull(bufRead, msg)
		if err != nil {
			return nil, true, err
			// 是io错误
		}
		if len(msg) == 0 ||
			msg[len(msg)-2] != '\r' ||
			msg[len(msg)-1] != '\n' ||
			len(msg)-2 != int(state.bulkLen) {
			return nil, false, errors.New("protocol error:" + string(msg))
		}
		state.bulkLen = 0
	}
	return msg, false, nil
}

// 读数组的
func parseMultiBulkHeader(msg []byte, state *readState) error {
	var err error
	var expectLine uint64
	expectLine, err = strconv.ParseUint(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		return errors.New("protocol error:" + string(msg))
	}
	if expectLine == 0 {
		state.expectedArgsCount = 0
		return nil
	} else if expectLine > 0 {
		state.msgType = msg[0]
		state.readingMultiLine = true
		state.expectedArgsCount = int(expectLine)
		// 可以存expectLine个[]byte
		state.args = make([][]byte, 0, expectLine)
		return nil
	} else {
		return errors.New("protocol error:" + string(msg))
	}
}

func parseBulkHeader(msg []byte, state *readState) error {
	var err error
	state.bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		return errors.New("protocol error:" + string(msg))
	}
	// 空字符串
	if state.bulkLen == -1 {
		state.expectedArgsCount = 0
		return nil
	} else if state.bulkLen > 0 {
		state.msgType = msg[0]
		state.expectedArgsCount = 1
		state.readingMultiLine = true
		state.args = make([][]byte, 0, 1)
		return nil
	} else {
		return errors.New("protocol error:" + string(msg))
	}
}

// 处理+，-，：正常错误数字
func parseSingleLineReply(msg []byte) (resp.Reply, error) {
	str := strings.TrimSuffix(string(msg), "\r\n")
	var result resp.Reply
	switch msg[0] {
	case '+':
		result = reply.NewStatusReply(str[1:])
	case '-':
		result = reply.NewErrReply(str[1:])
	case ':':
		value, err := strconv.ParseInt(str[1:], 10, 64)
		if err != nil {
			return nil, errors.New("protocol error:" + string(msg))
		}
		result = reply.NewIntReply(value)
	}
	return result, nil
}

// 读数组的后续字符串和字符串的内容
// *2\r\n$3\r\nset\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
// $3\r\nset\r\n
func readBody(msg []byte, state *readState) error {
	line := msg[0 : len(msg)-2]
	var err error
	if line[0] == '$' { // 是长度的标识
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return errors.New("protocol error:" + string(msg))
		}
		// $0\r\n
		if state.bulkLen <= 0 {
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else { // 读内容
		state.args = append(state.args, line)
	}
	return nil
}
