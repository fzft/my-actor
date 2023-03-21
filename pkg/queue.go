package pkg

import (
	"sync"
)

type Queue struct {
	items    []any
	cond     *sync.Cond
	capacity int
}

func NewQueue(cap int) *Queue {
	return &Queue{
		items: make([]any, 0, cap),
		cond:  sync.NewCond(&sync.Mutex{}),
	}
}

func (q *Queue) Enqueue(item any) {
	q.cond.L.Lock()
	q.items = append(q.items, item)
	q.cond.L.Unlock()

	q.cond.Signal()
}

func (q *Queue) Dequeue() any {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	for len(q.items) == 0 {
		q.cond.Wait()
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item
}
