package reply

/*
type UnknownErrorReply struct {
}

var unknownErrorReply = []byte("-Err unknown" + CRLF)

func (u UnknownErrorReply) Error() string {
	return "Err unknown"
}

func (u UnknownErrorReply) ToBytes() ([]byte, error) {
	return unknownErrorReply, nil
}

type ArgNumErrReply struct {
	cmd string
}

func (r *ArgNumErrReply) Error() string {
	return "ERR wrong number of arguments for '" + r.cmd + "' command"

}

func (r *ArgNumErrReply) ToBytes() ([]byte, error) {
	return []byte("-ERR wrong number of arguments for '" + r.cmd + "' command" + CRLF), nil
}

func NewArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{cmd: cmd}
}

// SyntaxErrReply represents meeting unexpected arguments
type SyntaxErrReply struct{}

var syntaxErrBytes = []byte("-Err syntax error" + CRLF)
var theSyntaxErrReply = &SyntaxErrReply{}

// NewSyntaxErrReply creates syntax error
func NewSyntaxErrReply() *SyntaxErrReply {
	return theSyntaxErrReply
}

// ToBytes marshals redis.Reply
func (r *SyntaxErrReply) ToBytes() ([]byte, error) {
	return syntaxErrBytes, nil
}

func (r *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

// WrongTypeErrReply represents operation against a key holding the wrong kind of value
type WrongTypeErrReply struct{}

var wrongTypeErrBytes = []byte("-WRONGTYPE Operation against a key holding the wrong kind of value" + CRLF)

// ToBytes marshals redis.Reply
func (r *WrongTypeErrReply) ToBytes() ([]byte, error) {
	return wrongTypeErrBytes, nil
}

func (r *WrongTypeErrReply) Error() string {
	return "WRONGTYPE Operation against a key holding the wrong kind of value"
}

// ProtocolErr

// ProtocolErrReply represents meeting unexpected byte during parse requests
type ProtocolErrReply struct {
	Msg string
}

// ToBytes marshals redis.Reply
func (r *ProtocolErrReply) ToBytes() ([]byte, error) {
	return []byte("-ERR Protocol error: " + r.Msg + CRLF), nil
}

func (r *ProtocolErrReply) Error() string {
	return "ERR Protocol error: CRLF" + r.Msg
}
*/
