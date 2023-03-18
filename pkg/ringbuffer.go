package pkg

import (
	"sync"
)

type RingBuffer[T any] struct {
	buffer   []T
	capacity int
	count    int
	head     int
	tail     int
	mu       sync.RWMutex
}

func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	return &RingBuffer[T]{
		buffer:   make([]T, capacity),
		capacity: capacity,
		count:    0,
		head:     0,
		tail:     0,
	}
}

func (rb *RingBuffer[T]) isEmpty() bool {
	return rb.count == 0
}

func (rb *RingBuffer[T]) IsEmpty() bool {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	return rb.isEmpty()
}

func (rb *RingBuffer[T]) isFull() bool {
	return rb.count == rb.capacity
}

func (rb *RingBuffer[T]) IsFull() bool {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	return rb.isFull()
}

func (rb *RingBuffer[T]) Enqueue(value T) bool {
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

func (rb *RingBuffer[T]) Dequeue() (T, bool) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.isEmpty() {
		var empty T
		return empty, false
	}

	value := rb.buffer[rb.head]
	rb.head = (rb.head + 1) % rb.capacity
	rb.count--

	return value, true
}

func (rb *RingBuffer[T]) Peek() (T, bool) {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	if rb.isEmpty() {
		var empty T
		return empty, false
	}

	value := rb.buffer[rb.head]
	return value, true
}
