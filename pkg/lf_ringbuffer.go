package pkg

import (
	"sync/atomic"
	"unsafe"
)

type Value struct {
	data    any
	version uint64
}

type LockFreeRingBuffer struct {
	buffer   []unsafe.Pointer
	head     uint64
	tail     uint64
	capacity uint64
}

func NewLockFreeRingBuffer(capacity int) *LockFreeRingBuffer {
	rb := &LockFreeRingBuffer{
		buffer:   make([]unsafe.Pointer, capacity),
		capacity: uint64(capacity),
	}
	for i := 0; i < capacity; i++ {
		rb.buffer[i] = unsafe.Pointer(&Value{})
	}
	return rb
}

func (rb *LockFreeRingBuffer) Enqueue(val any) bool {
	for {
		head := atomic.LoadUint64(&rb.head)
		tail := atomic.LoadUint64(&rb.tail)
		if (tail+1)%rb.capacity == head%rb.capacity {
			return false
		}

		expected := tail
		if atomic.CompareAndSwapUint64(&rb.tail, expected, (tail+1)%rb.capacity) {
			idx := tail % rb.capacity
			newValue := &Value{data: val, version: tail + rb.capacity}
			atomic.StorePointer(&rb.buffer[idx], unsafe.Pointer(newValue))
			return true
		}
	}
}

func (rb *LockFreeRingBuffer) Dequeue() (any, bool) {
	for {
		head := atomic.LoadUint64(&rb.head)
		tail := atomic.LoadUint64(&rb.tail)
		if head == tail {
			return nil, false
		}

		idx := head % rb.capacity
		value := (*Value)(atomic.LoadPointer(&rb.buffer[idx]))

		if value.version <= head {
			continue
		}

		expected := head
		if atomic.CompareAndSwapUint64(&rb.head, expected, (head+1)%rb.capacity) {
			return value.data, true
		}
	}
}

func (rb *LockFreeRingBuffer) IsEmpty() bool {
	head := atomic.LoadUint64(&rb.head)
	tail := atomic.LoadUint64(&rb.tail)
	return head == tail
}

func (rb *LockFreeRingBuffer) IsFull() bool {
	head := atomic.LoadUint64(&rb.head)
	tail := atomic.LoadUint64(&rb.tail)
	return (tail+1)%rb.capacity == head%rb.capacity
}

func (rb *LockFreeRingBuffer) Capacity() int {
	return int(rb.capacity) - 1
}
