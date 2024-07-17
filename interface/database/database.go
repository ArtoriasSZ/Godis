package database

import "Godis/interface/resp"

type CmdLine = [][]byte

type Database interface {
	Exec(client resp.Connection, args [][]byte) resp.Reply
	Close()
	AfterClientClose()
}

// DataEntity 指代各种数据类型
type DataEntity struct {
	Data interface{}
}
