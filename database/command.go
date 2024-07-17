package database

import "strings"

// 只在一开初始化，后面都是只读的，所以不用用sync.map
var cmdTable = map[string]*command{}

type command struct {
	executor ExecFunc
	arity    int
}

func RegisterCommand(name string, executor ExecFunc, arity int) {
	name = strings.ToLower(name)
	cmdTable[name] = &command{
		executor: executor,
		arity:    arity,
	}
}
