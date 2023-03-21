package main

// Define the Actor interface
type Actor interface {
	Receive(msg interface{})
}

// Define an actor struct that implements the Actor interface
type MyActor struct {
	state string
}

func (a *MyActor) Receive(msg any) {
}

func main() {
}
