// Package queue A linked list based Queue, use the NewQueue method to allocate a new queue
package queue

import (
	"sync"
	"time"
)

// INode type interface for generic Node,
// These types of nodes can be created using the generic declaration.
type INode interface {
	int | float32 | float64 | string
}

// Node Generic Doubly-Linked list.
type Node[T INode] struct {
	Key       *string
	Value     T
	CreatedAt time.Time
	TTL       *time.Duration
	Next      *Node[T]
	Prev      *Node[T]
}

// Queue Doubly-Linked list.
type Queue[T INode] struct {
	mu    sync.RWMutex
	Len   uint64
	Start *Node[T]
	Last  *Node[T]
}

// NewQueue Creates a New Queue with a given type.
func NewQueue[T INode]() *Queue[T] {
	return &Queue[T]{
		Len:   0,
		Start: nil,
		Last:  nil,
	}
}

func (q *Queue[T]) pushEnd(n *Node[T]) {
	if q.Start == nil {
		q.Start = n
		q.Last = n
		n.Prev = nil
		n.Next = nil
	} else {
		q.Last.Next = n
		n.Prev = q.Last
		q.Last = n
	}

	q.Len++
}

// PushEnd push a node to the end of queue.
func (q *Queue[T]) PushEnd(n *Node[T]) {
	n.CreatedAt = time.Now()
	q.mu.Lock()
	defer q.mu.Unlock()
	q.pushEnd(n)
}

func (q *Queue[T]) pushStart(n *Node[T]) {
	if q.Start == nil {
		q.Start = n
		q.Last = n
		n.Prev = nil
		n.Next = nil
	} else {
		q.Start.Prev = n
		n.Next = q.Start
		n.Prev = nil
		q.Start = n
	}

	q.Len++
}

// PushStart push a node to the start of the queue.
func (q *Queue[T]) PushStart(n *Node[T]) {
	n.CreatedAt = time.Now()
	q.mu.Lock()
	defer q.mu.Unlock()
	q.pushStart(n)
}

// Replaces the Start node values with new key and value
// and remove it from queue and add it to last so it becomes most recently used.
func (q *Queue[T]) LruMove(val T, key *string) {
	q.mu.RLock()
	if q.Len < 1 {
		q.mu.Unlock()
		return
	}

	q.mu.Lock()
	defer q.mu.Unlock()
	if q.Len == 1 {
		q.Start.Value = val
		q.Start.Key = key
	} else {
		n := q.Start
		n.Value = val
		n.Key = key

		// remove start and make the start->next as start
		q.Start = n.Next
		q.Start.Prev = nil

		// make removed start node n to end
		q.pushEnd(n)
	}
}
