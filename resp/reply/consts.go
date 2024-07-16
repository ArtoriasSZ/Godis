package reply

// pong回复
type pongReply struct {
}

var pongBytes = []byte("+PONG\r\n")

func (p pongReply) ToBytes() []byte {
	return pongBytes
}

func NewPongReply() *pongReply {
	return &pongReply{}
}

// ok回复
type okReply struct {
}

var okBytes = []byte("+OK\r\n")

func (p okReply) ToBytes() []byte {
	return okBytes
}
func NewOkReply() *okReply {
	return &okReply{}
}

// 空字符串回复
type NullBulkReply struct {
}

var nullBulkBytes = []byte("$-1\r\n")

func (p NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}
func NewNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

// 空数组
type EmptyMultiBulkReply struct {
}

var emptyMultBulkReply = []byte("*0\r\n")

func (p EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultBulkReply
}

func NewEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

// 空回复

type NoReply struct {
}

var noReply = []byte("")

func (p NoReply) ToBytes() []byte {
	return noReply
}
func NewNoReply() *NoReply {
	return &NoReply{}
}
