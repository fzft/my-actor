package pkg

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	cacheLineSize = 64
	paddingSize   = cacheLineSize - 8
)

type PaddedUint64 struct {
	value   uint64
	padding [paddingSize]byte
}

type WaitStrategy interface {
	Wait(sequence uint64, cursor *PaddedUint64) bool
}

type BusySpinWaitStrategy struct{}

func (w *BusySpinWaitStrategy) Wait(sequence uint64, cursor *PaddedUint64) bool {
	for {
		current := atomic.LoadUint64(&cursor.value)
		if current >= sequence {
			return true
		}
	}
}

type YieldingWaitStrategy struct{}

func (w *YieldingWaitStrategy) Wait(sequence uint64, cursor *PaddedUint64) bool {
	for {
		current := atomic.LoadUint64(&cursor.value)
		if current >= sequence {
			return true
		}
		time.Sleep(time.Microsecond)
	}
}

type EventHandler[T any] func(sequence uint64, value T)

type Disruptor[T any] struct {
	buffer       *LockFreeRingBuffer[T]
	cursor       PaddedUint64
	barrier      sync.WaitGroup
	eventHandler EventHandler[T]
}

func NewDisruptor[T any](bufferSize int, eventHandlers EventHandler[T]) *Disruptor[T] {
	return &Disruptor[T]{
		buffer:       NewLockFreeRingBuffer[T](bufferSize),
		eventHandler: eventHandlers,
	}
}

func (d *Disruptor[T]) Publish(event T) {
	atomic.AddUint64(&d.cursor.value, 1)
	d.buffer.Enqueue(event)
}

func (d *Disruptor[T]) Start() {
	d.barrier.Add(1)

	go func() {
		sequence := uint64(0)

		for {
			event, ok := d.buffer.Dequeue()
			if ok {
				d.eventHandler(sequence, *event)
				atomic.StoreUint64(&d.cursor.value, sequence)
			}
			sequence++
		}
	}()
}

func (d *Disruptor[T]) Stop() {
	d.barrier.Wait()
}
