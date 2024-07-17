package database

import (
	"Godis/config"
	"Godis/interface/resp"
	"Godis/lib/logger"
	"Godis/resp/reply"
	"strconv"
	"strings"
)

type Database struct {
	dbSet []*DB
}

func NewDatabase() *Database {
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	database := &Database{
		dbSet: make([]*DB, config.Properties.Databases),
	}
	for i := range database.dbSet {
		db := NewDB()
		db.index = i
		database.dbSet[i] = db
	}
	return database
}

// select
// get
// select
func (db *Database) Exec(client resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	if cmdName == "select" {
		if len(args) != 2 {
			return reply.NewErrReply("ERR wrong number of arguments for 'select' command")
		}
		return execSelect(client, *db, args[1:])
	}
	index := client.GetDBIndex()
	database := db.dbSet[index]
	return database.Exec(client, args)
}

func (db *Database) Close() {
	//TODO implement me
	panic("implement me")
}

func (db *Database) AfterClientClose() {
	//TODO implement me
	panic("implement me")
}

// select 1
func execSelect(c resp.Connection, database Database, args [][]byte) resp.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.NewErrReply("ERR invalid DB index")
	}
	if dbIndex >= len(database.dbSet) {
		return reply.NewErrReply("ERR DB index is out of range")
	}
	err = c.SelectDB(dbIndex)
	if err != nil {
		return reply.NewErrReply("ERR invalid DB index")
	}
	return reply.NewOkReply()
}
