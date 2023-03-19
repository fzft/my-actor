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

	// PreStart before the actor starts processing messages
	PreStart()

	// Receive is the method that defines the actor's behavior when it receives a message.
	Receive(ctx Context, msg any) error

	// PostStop after the actor stops processing messages
	PostStop()

	// Children returns the actor's children
	Children() []Actor
}

type DefaultActor struct {
}

func (d *DefaultActor) String() string {
	//TODO implement me
	panic("implement me")
}

func (d *DefaultActor) PreStart() {
	//TODO implement me
	panic("implement me")
}

func (d *DefaultActor) Receive(ctx Context, msg any) error {
	//TODO implement me
	panic("implement me")
}

func (d *DefaultActor) PostStop() {
	//TODO implement me
	panic("implement me")
}

func (d *DefaultActor) Children() []Actor {
	//TODO implement me
	panic("implement me")
}

func NewDefaultActor() *DefaultActor {
	return &DefaultActor{}
}
