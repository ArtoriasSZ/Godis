package dict

import "sync"

type SyncDict struct {
	m sync.Map
}

func NewSyncDict() *SyncDict {
	return &SyncDict{}
}

func (dict *SyncDict) Get(key string) (val interface{}, exists bool) {
	value, ok := dict.m.Load(key)
	return value, ok
}

func (dict *SyncDict) len() int {
	lenth := 0
	dict.m.Range(func(key, value interface{}) bool {
		lenth++
		return true
	})
	return lenth
}

func (dict *SyncDict) Put(key string, val interface{}) (result int) {
	_, exists := dict.m.Load(key)
	dict.m.Store(key, val)
	if exists {
		return 0
	}
	// 新的
	return 1
}

func (dict *SyncDict) PutIfAbsent(key string, val interface{}) (result int) {
	_, exists := dict.m.Load(key)
	if exists {
		return 0
	}
	// 新的
	dict.m.Store(key, val)
	return 1
}

func (dict *SyncDict) PutIfExists(key string, val interface{}) (result int) {
	_, exists := dict.m.Load(key)
	if exists {
		dict.m.Store(key, val)
		return 1
	}
	// 不存在
	return 0
}

func (dict *SyncDict) Remove(key string) (result int) {
	_, exists := dict.m.Load(key)
	dict.m.Delete(key)
	if exists {
		return 1
	}
	return 0
}

func (dict *SyncDict) ForEach(consumer Consumer) {
	dict.m.Range(func(key, value interface{}) bool {
		consumer(key.(string), value)
		return true
	})
}

func (dict *SyncDict) Keys() []string {
	keys := make([]string, dict.len())
	var index int
	dict.m.Range(func(key, value interface{}) bool {
		keys[index] = key.(string)
		index++
		return true
	})
	return keys
}

// 可以重复的
func (dict *SyncDict) RandomKeys(limit int) []string {
	keys := make([]string, limit)
	for i := 0; i < limit; i++ {
		dict.m.Range(func(key, value interface{}) bool {
			keys[i] = key.(string)
			i++
			return false
		})
	}
	return keys
}

func (dict *SyncDict) RandomDistinctKey(limit int) []string {
	keys := make([]string, limit)
	var index int
	dict.m.Range(func(key, value interface{}) bool {
		keys[index] = key.(string)
		index++
		if index == limit-1 {
			return false
		}
		return true
	})
	return keys
}

func (dict *SyncDict) Clear() {
	*dict = *NewSyncDict()
}
