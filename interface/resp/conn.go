package resp

type Connection interface {
	Writer([]byte) error
	GetDBindex() int
	SelectDB(index int) error
}
