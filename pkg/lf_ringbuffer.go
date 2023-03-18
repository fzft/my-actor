package pkg

import (
	"sync/atomic"
	"unsafe"
)

type Value[T any] struct {
	data    T
	version uint64
}

type LockFreeRingBuffer[T any] struct {
	buffer   []unsafe.Pointer
	head     uint64
	tail     uint64
	capacity uint64
}

func NewLockFreeRingBuffer[T any](capacity int) *LockFreeRingBuffer[T] {
	rb := &LockFreeRingBuffer[T]{
		buffer:   make([]unsafe.Pointer, capacity),
		capacity: uint64(capacity),
	}
	for i := 0; i < capacity; i++ {
		rb.buffer[i] = unsafe.Pointer(&Value[T]{})
	}
	return rb
}

func (rb *LockFreeRingBuffer[T]) Enqueue(val T) bool {
	for {
		head := atomic.LoadUint64(&rb.head)
		tail := atomic.LoadUint64(&rb.tail)
		if (tail+1)%rb.capacity == head%rb.capacity {
			return false
		}

		expected := tail
		if atomic.CompareAndSwapUint64(&rb.tail, expected, (tail+1)%rb.capacity) {
			idx := tail % rb.capacity
			newValue := &Value[T]{data: val, version: tail + rb.capacity}
			atomic.StorePointer(&rb.buffer[idx], unsafe.Pointer(newValue))
			return true
		}
	}
}

func (rb *LockFreeRingBuffer[T]) Dequeue() (*T, bool) {
	for {
		head := atomic.LoadUint64(&rb.head)
		tail := atomic.LoadUint64(&rb.tail)
		if head == tail {
			return nil, false
		}

		idx := head % rb.capacity
		value := (*Value[T])(atomic.LoadPointer(&rb.buffer[idx]))

		if value.version <= head {
			continue
		}

		expected := head
		if atomic.CompareAndSwapUint64(&rb.head, expected, (head+1)%rb.capacity) {
			return &value.data, true
		}
	}
}

func (rb *LockFreeRingBuffer[T]) IsEmpty() bool {
	head := atomic.LoadUint64(&rb.head)
	tail := atomic.LoadUint64(&rb.tail)
	return head == tail
}

func (rb *LockFreeRingBuffer[T]) IsFull() bool {
	head := atomic.LoadUint64(&rb.head)
	tail := atomic.LoadUint64(&rb.tail)
	return (tail+1)%rb.capacity == head%rb.capacity
}

func (rb *LockFreeRingBuffer[T]) Capacity() int {
	return int(rb.capacity) - 1
}
