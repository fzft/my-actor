package pkg

import (
	"fmt"
	"testing"
)

func TestNewDisruptor(t *testing.T) {
	handler := &SimpleHandler{}
	disruptor := NewDisruptor(1024, handler)

	disruptor.Start()

	for i := 0; i < 100; i++ {
		disruptor.Publish(Event{ID: i, Data: fmt.Sprintf("Event %d", i)})
	}

	disruptor.Stop()
}
