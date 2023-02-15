package queue_test

import (
	"testing"

	"github.com/dropdevrahul/gocache/gocache/queue"
	"github.com/stretchr/testify/assert"
)

func TestPushEnd(t *testing.T) {
	t.Run("queue-push-end", func(t *testing.T) {
		q := queue.NewQueue[string]()
		n := &queue.Node[string]{
			Value: "abc",
		}

		q.PushEnd(n)
		assert.Equal(t, q.Len, 1)
		assert.Equal(t, q.Start, n)
		assert.Equal(t, q.Last, n)

		n2 := &queue.Node[string]{
			Value: "def",
		}

		q.PushEnd(n2)
		assert.Equal(t, q.Len, 2)
		assert.Equal(t, q.Start, n)
		assert.Equal(t, q.Last, n2)
	})
}

func TestPushStart(t *testing.T) {
	t.Run("queue-push-end", func(t *testing.T) {
		q := queue.NewQueue[string]()
		q.PushEnd(&queue.Node[string]{
			Value: "abc",
		})
	})
}
