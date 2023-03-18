package pkg

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestDisruptor(t *testing.T) {
	bufferSize := 1024
	numEvents := 1000

	handler := func(sequence uint64, value int) {
		time.Sleep(10 * time.Millisecond)
		fmt.Printf("Processed value: %d\n", value)
	}

	disruptor := NewDisruptor(bufferSize, handler)

	disruptor.Start()
	defer disruptor.Stop()

	var wg sync.WaitGroup
	wg.Add(numEvents)

	for i := 0; i < numEvents; i++ {
		go func(i int) {
			defer wg.Done()
			disruptor.Publish(i)
		}(i)
	}

	wg.Wait()

}
