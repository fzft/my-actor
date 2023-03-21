package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLFRingBufferEnqueueDequeue(t *testing.T) {
	rb := NewLockFreeRingBuffer(5)

	rb.Enqueue(1)
	rb.Enqueue(2)
	rb.Enqueue(3)

	val, ok := rb.Dequeue()
	if !ok || val != 1 {
		t.Errorf("Dequeue failed, expected 1, got %v", val)
	}

	val, ok = rb.Dequeue()
	if !ok || val != 2 {
		t.Errorf("Dequeue failed, expected 2, got %v", val)
	}

	val, ok = rb.Dequeue()
	if !ok || val != 3 {
		t.Errorf("Dequeue failed, expected 3, got %v", val)
	}

	_, ok = rb.Dequeue()
	if ok {
		t.Errorf("Dequeue should have failed")
	}
}

func TestLFRingBufferIsEmpty(t *testing.T) {
	rb := NewLockFreeRingBuffer(3)

	if !rb.IsEmpty() {
		t.Error("Ring buffer should be empty")
	}

	rb.Enqueue(1)
	rb.Enqueue(2)

	if rb.IsEmpty() {
		t.Error("Ring buffer should not be empty")
	}

	val1, _ := rb.Dequeue()
	val2, _ := rb.Dequeue()
	assert.Equal(t, 1, val1)
	assert.Equal(t, 2, val2)

	if !rb.IsEmpty() {
		t.Error("Ring buffer should be empty")
	}
}

func TestLFRingBufferIsFull(t *testing.T) {
	rb := NewLockFreeRingBuffer(3)

	rb.Enqueue(1)
	rb.Enqueue(2)
	assert.True(t, rb.IsFull())

}

func TestLFRingBufferCapacity(t *testing.T) {
	rb := NewLockFreeRingBuffer(5)

	if rb.Capacity() != 4 {
		t.Errorf("Ring buffer capacity should be 4, got %d", rb.Capacity())
	}

	rb2 := NewLockFreeRingBuffer(10)

	if rb2.Capacity() != 9 {
		t.Errorf("Ring buffer capacity should be 9, got %d", rb2.Capacity())
	}
}

// TestLFRingBufferConcurrent tests concurrent access to the ring buffer
func TestLFRingBufferConcurrent(t *testing.T) {
	rb := NewLockFreeRingBuffer(5)

	for {
		go func() {
			rb.Enqueue(1)
		}()

		go func() {
			rb.Dequeue()
		}()
	}
}
