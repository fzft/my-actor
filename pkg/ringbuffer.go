package pkg

import (
	"sync"
)

type RingBuffer struct {
	buffer   []any
	capacity int
	count    int
	head     int
	tail     int
	mu       sync.RWMutex
}

func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		buffer:   make([]any, capacity),
		capacity: capacity,
		count:    0,
		head:     0,
		tail:     0,
	}
}

func (rb *RingBuffer) isEmpty() bool {
	return rb.count == 0
}

func (rb *RingBuffer) IsEmpty() bool {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	return rb.isEmpty()
}

func (rb *RingBuffer) isFull() bool {
	return rb.count == rb.capacity
}

func (rb *RingBuffer) IsFull() bool {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	return rb.isFull()
}

func (rb *RingBuffer) Enqueue(value any) bool {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.isFull() {
		return false
	}

	rb.buffer[rb.tail] = value
	rb.tail = (rb.tail + 1) % rb.capacity
	rb.count++

	return true
}

func (rb *RingBuffer) Dequeue() (any, bool) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.isEmpty() {
		var empty any
		return empty, false
	}

	value := rb.buffer[rb.head]
	rb.head = (rb.head + 1) % rb.capacity
	rb.count--

	return value, true
}

func (rb *RingBuffer) Peek() (any, bool) {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	if rb.isEmpty() {
		var empty any
		return empty, false
	}

	value := rb.buffer[rb.head]
	return value, true
}
