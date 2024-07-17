package resp

type Connection interface {
	Writer([]byte) error
	GetDBIndex() int
	SelectDB(index int) error
}
