package internel

import (
	"fmt"
	"github.com/fzft/my-actor/pkg"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

// Mailbox is the interface for the actor's mailbox
// MessageOrdering is the ordering of messages in the mailbox
// Throttling is the throttling of messages in the mailbox
// Author: fzft
const defaultThrottle = 1024

type Mailbox interface {
	// Source Post a message to the mailbox
	Source(msg any) error

	Consume() chan Message
}

type DefaultMailbox struct {
	logger   *zap.SugaredLogger
	throttle chan struct{}

	q          *pkg.Queue
	bufferSize int
	lastSent   time.Time
}

func NewDefaultMailbox(logger *zap.SugaredLogger) *DefaultMailbox {
	inbox := &DefaultMailbox{
		// throttle is the throttling of messages in the mailbox, the default is 1024
		throttle: make(chan struct{}, defaultThrottle),
		q:        pkg.NewQueue(defaultThrottle),
		logger:   logger,
	}
	return inbox
}

func (d *DefaultMailbox) Source(msg any) error {
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
func (d *DefaultMailbox) Consume() chan Message {
	c := make(chan Message, defaultThrottle)

	go func() {
		for {
			item := d.q.Dequeue()
			c <- WrapMsg(uuid.New().String(), item)
			<-d.throttle
		}
	}()

	return c
}
