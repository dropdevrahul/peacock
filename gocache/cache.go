// Package gocache A thread safe LRU based Cache
package gocache

import (
	"errors"
	"sync"
	"time"

	"github.com/dropdevrahul/gocache/gocache/queue"
)

var ErrEmptyValue = errors.New("empty value")

type CacheData struct {
	QueueNode *queue.Node[string]
}

// Cache uses LRU by default
// MaxCapacity is the maximum number of items that the underlying queue can contain
// q is a double linked list based queue which enables LRU removal of keys.
type Cache struct {
	q           *queue.Queue[string]
	cm          map[string]CacheData
	MaxCapacity uint64
	mu          sync.RWMutex
}

func NewCache(cap uint64) *Cache {
	c := &Cache{
		MaxCapacity: cap,
		cm:          map[string]CacheData{},
		q:           queue.NewQueue[string](),
	}

	return c
}

// Len Returns the current number of items in Cache.
func (c *Cache) Len() uint64 {
	return c.q.Len
}

// Set set a given key with given bytes data in the cache given the data can be encoded into
// utf-8 string.
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

	return nil
}

// Get fetches a given string from the db returns the string and whether the
// the string was found.
func (c *Cache) Get(key *string) (string, bool) {
	c.mu.RLock()
	n, ok := c.cm[*key]
	defer c.mu.RUnlock()

	if n.QueueNode == nil {
		return "", false
	}

	val := n.QueueNode.Value

	return val, ok
}

func (c *Cache) GetNode(key *string) (*queue.Node[string], bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	n, ok := c.cm[*key]

	if !ok {
		return nil, false
	}

	val := n.QueueNode

	return val, ok
}

func (c *Cache) SetTTL(key *string, ttl *time.Duration) int8 {
	n, ok := c.GetNode(key)

	if !ok {
		return 0
	}

	n.TTL = ttl
	return 1
}

func (c *Cache) GetTTL(key *string) time.Duration {
	n, ok := c.GetNode(key)

	if !ok {
		return -2
	}

	if n.TTL == nil {
		return -1
	}

	elapsed := time.Now().Sub(n.CreatedAt)
	ttl := *n.TTL - elapsed

	return ttl
}
