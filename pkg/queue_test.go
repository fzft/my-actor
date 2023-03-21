package pkg

import (
	"fmt"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	q := NewQueue(8)

	go func() {
		for i := 0; i < 10; i++ {
			q.Enqueue(i)
		}
	}()

	time.Sleep(2 * time.Second)

	for i := 0; i < 10; i++ {
		item := q.Dequeue()
		fmt.Println(item)
	}

}
