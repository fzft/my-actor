package actor

import (
	"fmt"
	"github.com/fzft/my-actor/pkg"
)

// DAG  directed acyclic graph (DAG) where actors are the nodes and messages are the edges
// Step 1: Create a DAG
// Step 2: Add actors to the DAG
// Step 3: Add messages to the Mailbox
// Step 4: Execute the actors in the DAG

// Engine is the actor engine
type Engine[T Actor] struct {

	// Nodes is the list of nodes in the engine
	*pkg.DAG[T]

	// Mailbox is the mailbox of the engine
	Mailbox Mailbox
}

func NewEngine() *Engine[Actor] {
	return &Engine[Actor]{}
}

// Spawn spawns a new actor
func (e *Engine[Actor]) Spawn(actor Actor) *pkg.Node[Actor] {
	return e.AddNode(actor)
}

// Execute executes the actors in the engine
func (e *Engine[Actor]) Execute() {
	visited := make(map[*pkg.Node[Actor]]bool)
	var visit func(node *pkg.Node[Actor])
	visit = func(node *pkg.Node[Actor]) {
		if visited[node] {
			return
		}
		visited[node] = true

		// TODO: execute the actor
		fmt.Println("Executing actor:", node.Value)

		neighbors := e.Neighbors(node)
		for _, neighbor := range neighbors {
			visit(neighbor)
		}
	}

	for _, node := range e.Nodes {
		if !visited[node] {
			visit(node)
		}
	}
}

