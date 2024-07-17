package database

import (
	"Godis/datastruct/dict"
	"Godis/interface/database"
	"Godis/interface/resp"
	"Godis/resp/reply"
	"fmt"
	"strings"
)

type CmdLine = [][]byte

type DB struct {
	index int
	data  dict.Dict
}

type ExecFunc func(db *DB, args [][]byte) resp.Reply

func NewDB() *DB {
	db := &DB{
		index: 0,
		data:  dict.NewSyncDict(),
	}
	return db
}

func (db *DB) Exec(c resp.Connection, cmdLine CmdLine) resp.Reply {
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmdFunc, ok := cmdTable[cmdName]
	if !ok {
		return reply.NewErrReply(fmt.Sprintf("ERR unknown command '%s'", cmdName))
	}
	// 校验参数个数
	if !validateArity(cmdFunc.arity, cmdLine) {
		return reply.NewArgNumErrReply(cmdName)
	}
	fun := cmdFunc.executor
	return fun(db, cmdLine[1:])
}

// Ser key var -> arity = 3
// Exists k1 k2... -> arity = -2(>=2)
func validateArity(arity int, cmdArgs [][]byte) bool {
	arityNum := len(cmdArgs)
	// todo: 不知道是否要减一
	if arity >= 0 {
		return arityNum == arity
	}
	return arityNum >= -arity
}

func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	row, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}
	return row.(*database.DataEntity), true
}

func (db *DB) PutEntity(key string, value *database.DataEntity) int {
	return db.data.Put(key, value)
}

func (db *DB) PutEntityIfExists(key string, value *database.DataEntity) int {
	return db.data.PutIfExists(key, value)
}
func (db *DB) PutEntityIfAbsent(key string, value *database.DataEntity) int {
	return db.data.PutIfAbsent(key, value)
}

func (db *DB) Remove(key string) {
	db.data.Remove(key)
}

func (db *DB) Removes(keys ...string) (deleted int) {
	deleted = 0
	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			db.Remove(key)
			deleted++
		}
	}
	return deleted
}

func (db *DB) Flush() {
	db.data.Clear()
}
