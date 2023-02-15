package gocache

import (
	"errors"
	"fmt"
	"sync"

	"github.com/dropdevrahul/gocache/gocache/queue"
)

var ErrEmptyValue = errors.New("empty value")

type CacheData struct {
	QueueNode *queue.Node[string]
}

// Cache uses LRU by default
type Cache struct {
	q           *queue.Queue[string]
	cm          map[string]CacheData
	MaxCapacity int
	mu          sync.RWMutex
}

func (c *Cache) Len() int {
	return c.q.Len
}

func (c *Cache) Set(key *string, data []byte) error {
	if len(data) == 0 {
		return ErrEmptyValue
	}

	dataS := string(data)

	c.mu.Lock()
	if c.q.Len < c.MaxCapacity {
		node := queue.Node[string]{
			Value: dataS,
			Key:   key,
		}
		c.q.PushEnd(&node)
	} else {
		rKey := c.q.Start
		if c.q.Start != nil {
			delete(c.cm, *rKey.Key)
		}
		c.q.LruMove(dataS, key)
	}

	c.cm[*key] = CacheData{
		QueueNode: c.q.Last,
	}

	defer c.mu.Unlock()

	fmt.Println("Set key: " + *key)

	return nil
}

func (c *Cache) Get(key *string) (string, bool) {
	c.mu.RLock()

	val, ok := c.cm[*key]

	defer c.mu.RUnlock()

	fmt.Println("Get key: " + *key)
	if val.QueueNode == nil {
		return "", ok
	}

	return val.QueueNode.Value, ok
}

func (c *Cache) Del(key *string) {
	c.mu.Lock()
	delete(c.cm, *key)
	defer c.mu.Unlock()
}
