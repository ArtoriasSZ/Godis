package reply

import (
	"Godis/interface/resp"
	"bytes"
	"strconv"
)

var (
	nullBulkReplyBytes = []byte("$-1")

	// CRLF is the line separator of redis serialization protocol
	CRLF = "\r\n"
)

/* ---- 字符串回复 ---- */

// bulk 大量，主体
type BulkReply struct {
	Arg []byte
}

// NewBulkReply creates  BulkReply
func NewBulkReply(arg []byte) *BulkReply {
	return &BulkReply{
		Arg: arg,
	}
}

// "moody"$5\r\nmoody\r\n""
func (r *BulkReply) ToBytes() ([]byte, error) {
	if len(r.Arg) == 0 {
		return nullBulkReplyBytes, nil
	}
	return []byte("$" + strconv.Itoa(len(r.Arg)) + CRLF + string(r.Arg) + CRLF), nil
}

/* ---- 字符串数组回复 ---- */

type MultiBulkReply struct {
	Args [][]byte
}

// NewMultiBulkReply creates MultiBulkReply
func NewMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{
		Args: args,
	}
}

// ToBytes marshal redis.Reply
func (r *MultiBulkReply) ToBytes() ([]byte, error) {
	argLen := len(r.Args)
	var buf bytes.Buffer
	_, err := buf.WriteString("*" + strconv.Itoa(argLen) + CRLF)
	// 缓冲区溢出报错
	if err != nil {
		return nil, err
	}
	for _, arg := range r.Args {
		if arg == nil {
			_, err = buf.WriteString("$-1" + CRLF)
			if err != nil {
				return nil, err
			}
		} else {
			_, err = buf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
			if err != nil {
				return nil, err
			}
		}
	}
	return buf.Bytes(), nil
}

/* ---- 状态回复 ---- */

type StatusReply struct {
	Status string
}

// NewStatusReply creates StatusReply
func NewStatusReply(status string) *StatusReply {
	return &StatusReply{
		Status: status,
	}
}

// ToBytes marshal redis.Reply
func (r *StatusReply) ToBytes() ([]byte, error) {
	return []byte("+" + r.Status + CRLF), nil
}

/* ---- 数字回复 ---- */

// IntReply stores an int64 number
type IntReply struct {
	Code int64
}

// NewIntReply creates int reply
func NewIntReply(code int64) *IntReply {
	return &IntReply{
		Code: code,
	}
}

// ToBytes marshal redis.Reply
func (r *IntReply) ToBytes() ([]byte, error) {
	return []byte(":" + strconv.FormatInt(r.Code, 10) + CRLF), nil
}

/* ---- 错误回复 ---- */

// ErrorReply is an error and redis.Reply
type ErrorReply interface {
	Error() string
	ToBytes() ([]byte, error)
}

// StandardErrReply represents handler error
type StandardErrReply struct {
	Status string
}

// ToBytes marshal redis.Reply
func (r *StandardErrReply) ToBytes() ([]byte, error) {
	return []byte("-" + r.Status + CRLF), nil
}

func (r *StandardErrReply) Error() string {
	return r.Status
}

// NewErrReply creates StandardErrReply
func NewErrReply(status string) *StandardErrReply {
	return &StandardErrReply{
		Status: status,
	}
}

// IsErrorReply returns true if the given reply is error
func IsErrorReply(reply resp.Reply) (bool, error) {
	toBytes, err := reply.ToBytes()
	if err != nil {
		return false, err
	}
	return toBytes[0] == '-', nil
}
