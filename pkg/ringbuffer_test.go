package pkg

import (
	"testing"
)

func TestNewRingBuffer(t *testing.T) {
	rb := NewRingBuffer[int](5)

	if rb.capacity != 5 {
		t.Errorf("Expected cap to be 5, but got %d", rb.capacity)
	}

	if rb.head != 0 {
		t.Errorf("Expected head to be 0, but got %d", rb.head)
	}

	if rb.tail != 0 {
		t.Errorf("Expected tail to be 0, but got %d", rb.tail)
	}
}

func TestIsEmpty(t *testing.T) {
	rb := NewRingBuffer[int](5)

	if !rb.IsEmpty() {
		t.Error("Expected buffer to be empty, but it's not")
	}
}

func TestIsFull(t *testing.T) {
	rb := NewRingBuffer[int](2)
	rb.Enqueue(1)
	rb.Enqueue(2)

	if !rb.IsFull() {
		t.Error("Expected buffer to be full, but it's not")
	}
}

func TestEnqueue(t *testing.T) {
	rb := NewRingBuffer[int](3)

	if !rb.Enqueue(1) {
		t.Error("Failed to enqueue value 1")
	}

	if !rb.Enqueue(2) {
		t.Error("Failed to enqueue value 2")
	}

	if !rb.Enqueue(3) {
		t.Error("Failed to enqueue value 3")
	}

	if rb.Enqueue(4) {
		t.Error("Enqueue should have failed as the buffer is full")
	}
}

func TestDequeue(t *testing.T) {
	rb := NewRingBuffer[int](3)

	rb.Enqueue(1)
	rb.Enqueue(2)

	value, ok := rb.Dequeue()
	if !ok || value != 1 {
		t.Errorf("Expected to dequeue value 1, but got %d", value)
	}

	value, ok = rb.Dequeue()
	if !ok || value != 2 {
		t.Errorf("Expected to dequeue value 2, but got %d", value)
	}

	_, ok = rb.Dequeue()
	if ok {
		t.Error("Dequeue should have failed as the buffer is empty")
	}
}

func TestPeek(t *testing.T) {
	rb := NewRingBuffer[int](3)

	rb.Enqueue(1)
	rb.Enqueue(2)

	value, ok := rb.Peek()
	if !ok || value != 1 {
		t.Errorf("Expected to peek value 1, but got %d", value)
	}

	rb.Dequeue()

	value, ok = rb.Peek()
	if !ok || value != 2 {
		t.Errorf("Expected to peek value 2, but got %d", value)
	}

	rb.Dequeue()

	_, ok = rb.Peek()
	if ok {
		t.Error("Peek should have failed as the buffer is empty")
	}
}
