package pkg

import "testing"

func BenchmarkLockFreeRingBuffer(b *testing.B) {
	rb := NewLockFreeRingBuffer(1024)
	stopCh := make(chan struct{})
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				rb.Enqueue(i)
			} else {
				rb.Dequeue(stopCh)
			}
			i++
		}
	})
}

func BenchmarkLockBasedRingBuffer(b *testing.B) {
	rb := NewRingBuffer(1024)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				rb.Enqueue(i)
			} else {
				rb.Dequeue()
			}
			i++
		}
	})
}
