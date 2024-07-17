package database

import (
	"Godis/interface/resp"
	"Godis/lib/wildcard"
	"Godis/resp/reply"
)

// DEL
func execDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	deleted := db.Removes(keys...)
	return reply.NewIntReply(int64(deleted))
}

// EXISTS
func execExists(db *DB, args [][]byte) resp.Reply {
	result := int64(0)
	for _, arg := range args {
		key := string(arg)
		_, exists := db.GetEntity(key)
		if exists {
			result++
		}
	}
	return reply.NewIntReply(result)
}

// KEYS

func execKeys(db *DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.data.ForEach(func(key string, val interface{}) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})
	return reply.NewMultiBulkReply(result)
}

// FLUSHDB
func execFlushDB(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	return reply.NewOkReply()
}

// TYPE
func execType(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, ok := db.GetEntity(key)
	if !ok {
		return reply.NewBulkReply([]byte("none"))
	}
	switch entity.Data.(type) {
	case []byte:
		return reply.NewStatusReply("string")
	}
	// todo
	return reply.UnknownErrorReply{}
}

// RENAME
func execRename(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])
	entity, ok := db.GetEntity(src)
	if !ok {
		return reply.NewErrReply("no such key")
	}
	db.PutEntity(dest, entity)
	db.Remove(src)
	return reply.NewOkReply()
}

// RENAMENX
// 判断key2是否存在
func execRenameNX(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])
	_, ok := db.GetEntity(dest)
	if ok {
		return reply.NewIntReply(0)
	}
	entity, ok := db.GetEntity(src)
	if !ok {
		return reply.NewErrReply("no such key")
	}
	db.PutEntity(dest, entity)
	db.Remove(src)
	return reply.NewIntReply(1)
}

func init() {
	RegisterCommand("DEL", execDel, -2)
	RegisterCommand("EXISTS", execExists, -2)
	RegisterCommand("Keys", execKeys, 2)
	RegisterCommand("FlushDB", execFlushDB, -1) // flushdb a b c
	RegisterCommand("TYPE", execType, 2)
	RegisterCommand("RENAME", execRename, 3)
	RegisterCommand("RENAMENX", execRenameNX, 3)

}
