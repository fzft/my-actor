package internel

// Actor is the interface that wraps the basic Actor methods.

// Actor Type
// 1. RootActor
//         read from mailbox and execute the message, then send the message to the child actor
// 2. IntermediateActor
//         read from parent actor and execute the message, then send the message to the child actor
// 3. LeafActor
//         read from parent actor and execute the message, and store the result in the sinkPool

// Actor Receive: This method is used to define the actor's behavior when it receives a message.
// It takes a message object as its argument and returns an error object
type Actor interface {
	// String ActorId returns the actor's id
	String() string

	// Receive is the method that defines the actor's behavior when it receives a message.
	Receive(ctx *Context, msg any) (any, error)
}

// ErrHandlerActor is the interface that wraps the basic ErrHandlerActor methods.
type ErrHandlerActor interface {
	Actor

	// ErrHandler is the method that defines the actor's behavior when it receives a error.
	ErrHandler(ctx *Context, err error)
}

type PreStartHookActor interface {
	Actor
	// PreStart before the actor starts processing messages
	PreStart()
}

type PostStopHookActor interface {
	Actor

	// PostStop after the actor stops processing messages
	PostStop()
}

type PreHandleMsgHookActor interface {
	Actor

	// PreHandleMsg before the actor handles a message
	PreHandleMsg(ctx *Context, inMsg any)
}

type PostHandleMsgHookActor interface {
	Actor

	// PostHandleMsg after the actor handles a message
	PostHandleMsg(ctx *Context, outMsg any)
}

type DefaultActor struct {
}

func NewDefaultActor(name string) *DefaultActor {
	return &DefaultActor{}
}

func (d *DefaultActor) String() string {
	return "DefaultActor"
}

// PreStart is Hook that is called before the actor starts processing messages.
func (d *DefaultActor) PreStart() {
	//TODO implement me
	panic("implement me")
}

func (d *DefaultActor) Receive(ctx Context, msg any) error {
	//TODO implement me
	panic("implement me")
}

// PostStop is Hook that is called after the actor stops processing messages.
func (d *DefaultActor) PostStop() {
	//TODO implement me
	panic("implement me")
}
