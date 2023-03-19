package internel

// Suber is the interface that wraps the basic Subscribe methods.
// subscribe to a or more parent actor
type Suber interface {
	Sub(<-chan any)
}

// Puber is the interface that wraps the basic Publish methods.
// publish to a or more child actor
type Puber interface {
	Pub() <-chan any
}

// Context is the interface that wraps the basic Context methods.
// every actor has a context
// the context is used to communicate with the actor's parent and children
type Context struct {
	sub   Suber
	pub   Puber
	store Storer

	pid Pid
}

// NewContext returns a new Context
func NewContext(pid Pid) *Context {
	return &Context{pid: pid, store: NewMemoryStore()}
}
