package pkg

import (
	"sync"
)

type Queue[T any] struct {
	items    []T
	cond     *sync.Cond
	capacity int
}

func NewQueue[T any](cap int) *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0, cap),
		cond:  sync.NewCond(&sync.Mutex{}),
	}
}

func (q *Queue[T]) Enqueue(item T) {
	q.cond.L.Lock()
	q.items = append(q.items, item)
	q.cond.L.Unlock()

	q.cond.Signal()
}

func (q *Queue[T]) Dequeue() T {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	for len(q.items) == 0 {
		q.cond.Wait()
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item
}
