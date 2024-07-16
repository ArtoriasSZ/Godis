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

type PayLoad struct {
	Date resp.Reply // 客户端发送的和服务端发送的都是reply
	Err  error
}

type readState struct {
	readingMultiLine  bool     // 正在解析多行数据还是单行
	expectedArgsCount int      // 期待解析的参数个数
	msgType           byte     // 用户消息类型
	args              [][]byte // 用户传参
	bulkLen           int64    // 指令长度
}

func (r readState) finished() bool {
	return len(r.args) > 0 && len(r.args) == r.expectedArgsCount
}

func ParserStream(reader io.Reader) <-chan *PayLoad {
	ch := make(chan *PayLoad)
	go parser0(reader, ch)
	return ch
}

func parser0(reader io.Reader, ch chan *PayLoad) {
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
		if err != nil {
			if ioErr {
				ch <- &PayLoad{
					Err: err,
				}
				close(ch)
			}
			ch <- &PayLoad{
				Err: err,
			}
			state = readState{}
			continue
		}
		// 多行解析模式
		if !state.readingMultiLine {
			if msg[0] == '*' {
				err = parserMultiBulkHeader(msg, &state)
				if err != nil {
					ch <- &PayLoad{
						Err: err,
					}
					state = readState{}
					continue
				}
				if state.expectedArgsCount == 0 {
					ch <- &PayLoad{
						Date: reply.NewEmptyMultiBulkReply(),
					}
					state = readState{}
					continue
				}
			} else if msg[0] == '$' {
				err = parserBulkHeader(msg, &state)
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
			} else {
				var lineReply resp.Reply
				lineReply, err = parserSingleLineReply(msg)
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
			if !state.finished() {
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

func readline(bufRead *bufio.Reader, state *readState) ([]byte, bool, error) {
	// 1.没有读几个的预设，\r\n分割
	var msg []byte
	var err error
	if state.bulkLen == 0 {
		// 没有预设
		msg, err = bufRead.ReadBytes('\n')
		if err != nil {
			return nil, true, err
			// 是io错误
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol error:" + string(msg))
		}
	} else { // 2.前面有限制，明确要读多少个。不受\r\n预设
		msg = make([]byte, state.bulkLen+2)
		_, err = io.ReadFull(bufRead, msg)
		if err != nil {
			return nil, true, err
			// 是io错误
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' || msg[len(msg)-1] != '\n' {
			return nil, false, errors.New("protocol error:" + string(msg))
		}
		state.bulkLen = 0
	}
	return msg, false, nil
}

func parserMultiBulkHeader(msg []byte, state *readState) error {
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
		state.expectedArgsCount = int(expectLine)
		state.readingMultiLine = true
		state.args = make([][]byte, 0, expectLine)
		return nil
	} else {
		return errors.New("protocol error:" + string(msg))
	}
}

func parserBulkHeader(msg []byte, state *readState) error {
	var err error
	state.bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		return errors.New("protocol error:" + string(msg))
	}
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

func parserSingleLineReply(msg []byte) (resp.Reply, error) {
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

func readBody(msg []byte, state *readState) error {
	line := msg[0 : len(msg)-2]
	var err error
	if line[0] == '$' {
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return errors.New("protocol error:" + string(msg))
		}
		// $0\r\n
		if state.bulkLen <= 0 {
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		state.args = append(state.args, line)
	}
	return nil
}
