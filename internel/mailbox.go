package internel

import (
	"fmt"
	"github.com/fzft/my-actor/pkg"
	"time"
)

// Mailbox is the interface for the actor's mailbox
// MessageOrdering is the ordering of messages in the mailbox
// Throttling is the throttling of messages in the mailbox
// Author: fzft
const defaultThrottle = 1024

type Mailbox[T any] interface {
	// Source Post a message to the mailbox
	Source(msg T) error

	Consume() <-chan T
}

type DefaultMailbox[T any] struct {
	throttle chan struct{}

	q          *pkg.Queue[T]
	bufferSize int
	lastSent   time.Time
}

func NewDefaultMailbox[T any]() *DefaultMailbox[T] {
	inbox := &DefaultMailbox[T]{
		// throttle is the throttling of messages in the mailbox, the default is 1024
		throttle: make(chan struct{}, defaultThrottle),
		q:        pkg.NewQueue[T](defaultThrottle),
	}
	return inbox
}

func (d *DefaultMailbox[T]) Source(msg T) error {
	select {
	case d.throttle <- struct{}{}:
		d.lastSent = time.Now()
		d.q.Enqueue(msg)
	default:
		// Throttle channel is full, drop message
		return fmt.Errorf("throttle channel is full, drop message")
	}
	return nil
}

// Consume ...
func (d *DefaultMailbox[T]) Consume() <-chan T {
	c := make(chan T, defaultThrottle)

	go func() {
		for {
			item := d.q.Dequeue()
			c <- item
			<-d.throttle
		}
	}()

	return c
}
