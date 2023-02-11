package gocache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrEmptyValue = errors.New("empty value")

type CacheData struct {
	LastUsed  int64
	BytesData []byte
}

type Cache struct {
	cm map[string]CacheData
	mu sync.Mutex
}

func (c *Cache) Set(key *string, data []byte) error {
	if len(data) == 0 {
		return ErrEmptyValue
	}

	c.mu.Lock()
	c.cm[*key] = CacheData{
		LastUsed:  (time.Now()).UnixNano(),
		BytesData: data,
	}
	fmt.Println("Set key: " + *key)

	defer c.mu.Unlock()

	return nil
}

func (c *Cache) Get(key *string) (CacheData, bool) {
	c.mu.Lock()
	val, ok := c.cm[*key]
	defer c.mu.Unlock()
	fmt.Println("Get key: " + *key)
	return val, ok
}

func (c *Cache) Del(key *string) {
	c.mu.Lock()
	delete(c.cm, *key)
	defer c.mu.Unlock()
}

var cm = map[string]CacheData{}
var HashMapCache = &Cache{
	cm: map[string]CacheData{},
}
