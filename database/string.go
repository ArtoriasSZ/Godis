package database

import (
	"Godis/interface/database"
	"Godis/interface/resp"
	"Godis/resp/reply"
)

// GET
func execGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.NewNullBulkReply()
	}
	// 类型多的时候是可能失败的
	bytes, ok := entity.Data.([]byte)
	if !ok {
	}
	return reply.NewBulkReply(bytes)
}

// SET
func execSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := string(args[1])
	entity := &database.DataEntity{
		Data: []byte(value),
	}
	db.PutEntity(key, entity)
	return reply.NewOkReply()
}

// SETNX
func execSetNX(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := string(args[1])
	entity := &database.DataEntity{
		Data: []byte(value),
	}
	res := db.PutEntityIfAbsent(key, entity)
	return reply.NewIntReply(int64(res))
}

// GETSET
func execGetSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := args[1]
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.NewNullBulkReply()
	}
	// 类型多的时候是可能失败的
	db.PutEntity(key, &database.DataEntity{
		Data: value,
	})
	bytes, ok := entity.Data.([]byte)
	if !ok {
	}
	return reply.NewBulkReply(bytes)
}

// STRLEN
func execStrLen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.NewNullBulkReply()
	}
	// 类型多的时候是可能失败的
	bytes, ok := entity.Data.([]byte)
	if !ok {
	}
	return reply.NewIntReply(int64(len(bytes)))
}

func init() {
	RegisterCommand("get", execGet, 2)       // get k1
	RegisterCommand("set", execSet, 3)       // set k1 v1
	RegisterCommand("setnx", execSetNX, 3)   // setnx k1 v1
	RegisterCommand("getset", execGetSet, 3) // getset k1 v1
	RegisterCommand("strlen", execStrLen, 2) // strlen k1
}
