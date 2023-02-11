package queue

import "sync"

type INode interface {
	int | float32 | float64 | string
}

type Node[T INode] struct {
	Key   *string
	Value T
	Next  *Node[T]
	Prev  *Node[T]
}

type Queue[T INode] struct {
	mu    sync.Mutex
	Len   int
	Start *Node[T]
	Last  *Node[T]
}

func NewQueue[T INode]() *Queue[T] {
	return &Queue[T]{
		Len:   0,
		Start: nil,
		Last:  nil,
	}
}

func (q *Queue[T]) PushEnd(n *Node[T]) {
	q.mu.Lock()
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

	q.Len += 1
	q.mu.Unlock()
}

func (q *Queue[T]) PushStart(n *Node[T]) {
	q.mu.Lock()
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
	q.mu.Unlock()
}

func (q *Queue[T]) LruMove(val T, key *string) {

	q.mu.Lock()
	if q.Len < 1 {
		q.mu.Unlock()
		return
	}

	if q.Len == 1 {
		q.Start.Value = val
		q.Start.Key = key
	} else {

		n := q.Start
		n.Value = val
		n.Key = key
		q.Start = n.Next
		q.Start.Prev = nil
		q.Last.Next = n
		n.Prev = q.Last
		q.Last = n
	}

	q.mu.Unlock()
}
