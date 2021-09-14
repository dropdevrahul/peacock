package gocache

import (
	"sync"
	"time"
)

type CacheData struct {
    lastUsed int64
    bytesData []byte
    mu sync.Mutex
}


type Cache struct {
    cm map[string]CacheData
}


func (c* Cache) Set(key* string, data []byte) {
    if val, ok := c.cm[*key]; ok {
        val.mu.Lock()
        c.cm[*key] = CacheData{
            lastUsed: (time.Now()).UnixNano(),
            bytesData: data,
        }
        val.mu.Unlock()
        return
    }
    c.cm[*key] = CacheData{
        lastUsed: (time.Now()).UnixNano(),
        bytesData: data,
    }
    return
}

func (c* Cache) Get(key* string) (CacheData, bool) {
    val, ok := c.cm[*key]
    if ok {
        val.mu.Lock()
        defer val.mu.Unlock()
        return val, ok
    }
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
