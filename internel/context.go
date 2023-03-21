package internel

import (
	"github.com/fzft/my-actor/pkg"
	"go.uber.org/zap"
)

// InBox maintains a lock-free ring buffer for incoming messages
type InBox struct {
	buffer *pkg.LockFreeRingBuffer
}

func NewInBox(bufferSize int) *InBox {
	return &InBox{buffer: pkg.NewLockFreeRingBuffer(bufferSize)}
}

// Enqueue adds a message to the actor's inbox
func (i *InBox) Enqueue(msg any) {
	i.buffer.Enqueue(msg)
}

// Dequeue removes a message from the actor's inbox
func (i *InBox) Dequeue() (any, bool) {
	return i.buffer.Dequeue()
}

// Context is the interface that wraps the basic Context methods.
// every actor has a context
// the context is used to communicate with the actor's parent and children
type Context struct {
	logger *zap.SugaredLogger

	// Suber is the channel that the actor's parent uses to send messages to the actor
	Suber chan Message
	store Storer

	inbox InBox // the actor's inbox

	// pid is the actor's pid
	pid      string
	children []*Pid
}

// NewContext returns a new Context
func NewContext(logger *zap.SugaredLogger, pid string) *Context {
	ctx := &Context{pid: pid, store: NewMemoryStore(), inbox: *NewInBox(1024), logger: logger, Suber: make(chan Message)}
	go ctx.buffered()
	return ctx
}

// AddChild adds a child to the actor's children
func (c *Context) addChild(pid *Pid) {
	c.children = append(c.children, pid)
}

// Children returns the actor's children
func (c *Context) childActors() []*Pid {
	return c.children
}

// buffered buffer the incoming message from suber into the actor's inbox
func (c *Context) buffered() {
	for {
		select {
		case msg, ok := <-c.Suber:
			if ok {
				c.logger.Debugf("[%s] buffered %+v", c.pid, msg)
				c.inbox.Enqueue(msg)
			}
		}
	}
}

// setMailbox sets the root actor's mailbox
func (c *Context) setMailbox(mailbox Mailbox) {
	c.Suber = mailbox.Consume()
}

// broadcast the outgoing message from the actor's inbox to the puber
func (c *Context) broadcast(msg Message) {
	for _, child := range c.children {
		go func(child *Pid) {
			for {
				select {
				case child.context.Suber <- msg:
					c.logger.Debugf("[%s] broadcast %v -> [%s] ", c.pid, msg, child.context.pid)
				}
			}
		}(child)
	}
}
