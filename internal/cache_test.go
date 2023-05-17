package gocache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_SetKey(t *testing.T) {
	t.Parallel()
	cache := NewCache(10)

	key := "abc"
	err := cache.Set(&key, []byte("hello world"))

	assert.Equal(t, nil, err)
	assert.Equal(t, cache.Len(), uint64(1))
	assert.Equal(t, "hello world", cache.q.Start.Value)
	assert.Equal(t, (*time.Duration)(nil), cache.q.Start.TTL)
	assert.Equal(t, &key, cache.q.Start.Key)
	assert.Equal(t, cache.q.Start, cache.cm["abc"].QueueNode)
}

func Test_GetKey(t *testing.T) {
	t.Parallel()
	cache := NewCache(10)

	key := "abc"
	err := cache.Set(&key, []byte("hello world"))

	assert.Equal(t, nil, err)

	val, ok := cache.Get(&key)

	assert.Equal(t, true, ok)
	assert.Equal(t, "hello world", val)
}

func Test_GetKey_empty(t *testing.T) {
	t.Parallel()
	cache := NewCache(10)
	key := "abc"

	val, ok := cache.Get(&key)
	assert.Equal(t, false, ok)
	assert.Equal(t, "", val)
}

func Test_SetTTL_NoKey(t *testing.T) {
	t.Parallel()
	cache := NewCache(10)
	key := "abc"

	ttl := time.Second * 100
	r := cache.SetTTL(&key, &ttl)

	assert.Equal(t, int8(0), r)
}

func Test_SetTTL(t *testing.T) {
	t.Parallel()
	cache := NewCache(10)
	key := "abc"
	ttl := time.Second * 100

	cache.Set(&key, []byte("hello world"))

	r := cache.SetTTL(&key, &ttl)

	assert.Equal(t, int8(1), r)
	assert.Equal(t, ttl, *cache.q.Start.TTL)
}

func Test_GetTTL_NoKey(t *testing.T) {
	t.Parallel()
	cache := NewCache(10)
	key := "abc"

	r := cache.GetTTL(&key)

	assert.Equal(t, int8(-2), r)
}

func Test_GetTTL(t *testing.T) {
	cache := NewCache(10)
	key := "abc"
	ttl := time.Second * 100

	cache.Set(&key, []byte("hello world"))

	r := cache.GetTTL(&key)

	assert.Equal(t, time.Duration(-1), r)
	assert.Equal(t, (*time.Duration)(nil), cache.q.Start.TTL)

	cache.q.Start.TTL = &ttl

	time.Sleep(time.Second * 2)
	r = cache.GetTTL(&key)

	assert.Equal(t, (time.Second * 97).Truncate(time.Second), r.Truncate(time.Second))
}
