package pkg

import (
	"fmt"
	"sync"
	"time"
)

type Event struct {
	ID   int
	Data string
}

type EventHandler interface {
	HandleEvent(event Event)
}

type SimpleHandler struct{}

func (sh *SimpleHandler) HandleEvent(event Event) {
	fmt.Printf("Event ID: %d, Data: %s\n", event.ID, event.Data)
}

type Disruptor struct {
	rb        *RingBuffer[Event]
	handler   EventHandler
	waitGroup sync.WaitGroup
	stop      chan struct{}
}

func NewDisruptor(bufferSize int, handler EventHandler) *Disruptor {
	return &Disruptor{
		rb:        NewRingBuffer[Event](bufferSize),
		handler:   handler,
		waitGroup: sync.WaitGroup{},
		stop:      make(chan struct{}),
	}
}

func (d *Disruptor) Publish(event Event) {
	for !d.rb.Enqueue(event) {
		time.Sleep(1 * time.Microsecond)
	}
	d.waitGroup.Add(1)
}

func (d *Disruptor) Start() {
	go func() {
		for {
			select {
			case <-d.stop:
				return
			default:
				event, ok := d.rb.Dequeue()
				if ok {
					d.handler.HandleEvent(event)
					d.waitGroup.Done()
				}
			}
		}
	}()
}

func (d *Disruptor) Stop() {
	d.waitGroup.Wait()
	close(d.stop)
}
