package actor

// Actor is the interface that wraps the basic Actor methods.

// Actor Receive: This method is used to define the actor's behavior when it receives a message.
//It takes a message object as its argument and returns an error object

type Actor interface {
	// String ActorId returns the actor's id
	String() string

	// PreStart before the actor starts processing messages
	PreStart()

	// Receive is the method that defines the actor's behavior when it receives a message.
	Receive()

	// PostStop after the actor stops processing messages
	PostStop()

	// Supervise is used to define the actor's behavior when it encounters an error or exception
	Supervise()

	// Children returns the actor's children
	Children() []Actor

	Context() Context
}
