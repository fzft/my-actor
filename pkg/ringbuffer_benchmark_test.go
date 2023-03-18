package pkg

import "testing"

func BenchmarkLockFreeRingBuffer(b *testing.B) {
	rb := NewLockFreeRingBuffer[int](1024)
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

func BenchmarkLockBasedRingBuffer(b *testing.B) {
	rb := NewRingBuffer[int](1024)
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
