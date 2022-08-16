package gocache

import (
	"fmt"
	"sync"
	"time"
)

type CacheData struct {
	lastUsed  int64
	bytesData []byte
}

type Cache struct {
	cm map[string]CacheData
	mu sync.Mutex
}

func (c *Cache) Set(key *string, data []byte) {
	c.mu.Lock()
	c.cm[*key] = CacheData{
		lastUsed:  (time.Now()).UnixNano(),
		bytesData: data,
	}
	fmt.Println(*key)
	defer c.mu.Unlock()
	return
}

func (c *Cache) Get(key *string) (CacheData, bool) {
	c.mu.Lock()
	val, ok := c.cm[*key]
	defer c.mu.Unlock()
	return val, ok
}

//func (c* Cache) Del(key *string) {
//c.mu.lock()
//c.mu.unlock()
//}

var cm = map[string]CacheData{}
var HashMapCache = &Cache{
	cm: map[string]CacheData{},
}
