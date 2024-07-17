package database

import (
	"Godis/interface/resp"
	"Godis/resp/reply"
)

func Ping(db *DB, args [][]byte) resp.Reply {
	return reply.NewPongReply()
}

func init() {
	RegisterCommand("ping", Ping, 1)
}
