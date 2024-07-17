package database

import (
	"Godis/interface/resp"
	"Godis/resp/reply"
)

type echodatabase struct {
}

func NewEchodatabase() *echodatabase {
	return &echodatabase{}
}

func (e echodatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	return reply.NewMultiBulkReply(args)
}

func (e echodatabase) Close() {

}

func (e echodatabase) AfterClientClose() {

}
