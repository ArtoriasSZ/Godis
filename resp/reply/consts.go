package reply

// pong回复
type PongReply struct {
}

var pongBytes = []byte("+PONG" + CRLF)

func (p PongReply) ToBytes() ([]byte, error) {
	return pongBytes, nil
}

func NewPongReply() *PongReply {
	return &PongReply{}
}

// ok回复
type OkReply struct {
}

var okBytes = []byte("+OK" + CRLF)

func (p OkReply) ToBytes() ([]byte, error) {
	return okBytes, nil
}
func NewOkReply() *OkReply {
	return &OkReply{}
}

// 空字符串回复
type NullBulkReply struct {
}

var nullBulkBytes = []byte("$-1" + CRLF)

func (p NullBulkReply) ToBytes() ([]byte, error) {
	return nullBulkBytes, nil
}
func NewNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

// 空数组
type EmptyMultiBulkReply struct {
}

var emptyMultiBulkReply = []byte("*0" + CRLF)

func (p EmptyMultiBulkReply) ToBytes() ([]byte, error) {
	return emptyMultiBulkReply, nil
}

func NewEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

// 空回复

type NoReply struct {
}

var noReply = []byte("")

func (p NoReply) ToBytes() ([]byte, error) {
	return noReply, nil
}
func NewNoReply() *NoReply {
	return &NoReply{}
}
