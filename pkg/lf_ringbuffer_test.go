package pkg

import "testing"

func TestLFRingBufferEnqueueDequeue(t *testing.T) {
	rb := NewLockFreeRingBuffer[int](5)

	rb.Enqueue(1)
	rb.Enqueue(2)
	rb.Enqueue(3)

	val, ok := rb.Dequeue()
	if !ok || *val != 1 {
		t.Errorf("Dequeue failed, expected 1, got %v", val)
	}

	val, ok = rb.Dequeue()
	if !ok || *val != 2 {
		t.Errorf("Dequeue failed, expected 2, got %v", val)
	}

	val, ok = rb.Dequeue()
	if !ok || *val != 3 {
		t.Errorf("Dequeue failed, expected 3, got %v", val)
	}

	_, ok = rb.Dequeue()
	if ok {
		t.Errorf("Dequeue should have failed")
	}
}

func TestLFRingBufferIsEmpty(t *testing.T) {
	rb := NewLockFreeRingBuffer[int](3)

	if !rb.IsEmpty() {
		t.Error("Ring buffer should be empty")
	}

	rb.Enqueue(1)
	rb.Enqueue(2)

	if rb.IsEmpty() {
		t.Error("Ring buffer should not be empty")
	}

	rb.Dequeue()
	rb.Dequeue()

	if !rb.IsEmpty() {
		t.Error("Ring buffer should be empty")
	}
}

func TestLFRingBufferIsFull(t *testing.T) {
	rb := NewLockFreeRingBuffer[int](3)

	rb.Enqueue(1)
	rb.Enqueue(2)

	if rb.IsFull() {
		t.Error("Ring buffer should not be full")
	}

	rb.Enqueue(3)

	if !rb.IsFull() {
		t.Error("Ring buffer should be full")
	}
}

func TestLFRingBufferCapacity(t *testing.T) {
	rb := NewLockFreeRingBuffer[int](5)

	if rb.Capacity() != 4 {
		t.Errorf("Ring buffer capacity should be 4, got %d", rb.Capacity())
	}

	rb2 := NewLockFreeRingBuffer[int](10)

	if rb2.Capacity() != 9 {
		t.Errorf("Ring buffer capacity should be 9, got %d", rb2.Capacity())
	}
}
