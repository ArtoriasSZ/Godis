package resp

type Reply interface {
	// ToBytes
	// 与其在处理MultiBulkReply的时候因为忘记处理缓存区可能溢出的情况
	// 不如每个ToBytes都去判断error的情况
	ToBytes() ([]byte, error)
}
