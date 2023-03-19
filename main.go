package main

import (
	"fmt"
	"github.com/fzft/my-actor/internel"
)

// Define the Actor interface
type Actor interface {
	Receive(msg interface{})
}

// Define an actor struct that implements the Actor interface
type MyActor struct {
	state string
}

func (a *MyActor) Receive(msg any) {
	// Process the message
	switch msg := msg.(type) {
	case internel.Message:
		fmt.Printf("Received message: %s\n", msg.content)
	}
}

func main() {
	// Create a new actor
	actor := &MyActor{state: "Initial state"}

	// Send a message to the actor
	msg := Message{content: "New state"}
	actor.Receive(msg)

	// Print the actor's state
	fmt.Printf("Actor state: %s\n", actor.state)
}
