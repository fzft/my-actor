package internel

import (
	"fmt"
	"github.com/fzft/my-actor/pkg"
	"go.uber.org/zap"
)

// DAG  directed acyclic graph (DAG) where actors are the nodes
// DAG is allowed to have only one root actor
// every actor can have multiple parents and children
// actor is driven by messages
// between actors, there are edges that connect them
// the edges are unidirectional
// edges are used to send messages between actors
// use ring buffer to store messages between actors
// if received message is not handled, will discard it

// Step 1: Create a DAG
// Step 2: Add actors to the DAG
// Step 3: Add Edge between actors
// Step 3: Ready the actors in the DAG
// Step 4: Send messages to the DAG

// Engine is the actor engine
type Engine[T Actor] struct {
	logger *zap.SugaredLogger

	// Nodes is the list of nodes in the engine
	*pkg.DAG[T]

	// mailbox is the mailbox of the engine
	mailbox Mailbox[any]

	// nodeMaps is the map of actors
	nodeMaps map[string]*pkg.Node[T]

	// root is the root actor
	root T

	// sinkPool is store the result of the leaf actor
	sinkPool *SinkPool
}

func NewEngine() *Engine[Actor] {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.TimeKey = ""
	logger, _ := loggerConfig.Build()
	return &Engine[Actor]{
		logger:   logger.Sugar(),
		DAG:      pkg.NewDAG[Actor](),
		mailbox:  NewDefaultMailbox[any](),
		sinkPool: NewSinkPool(),
	}
}

// Spawn spawns a new actor
func (e *Engine[Actor]) Spawn(actor Actor) Pid {
	pid := NewPid(actor.String())
	node := e.AddNode(actor)
	e.nodeMaps[pid.uuid] = node
	return pid
}

func (e *Engine[Actor]) AddEdge(from, to Pid) error {
	fromNode, ok := e.nodeMaps[from.uuid]
	if !ok {
		return fmt.Errorf("from actor not found")
	}

	toNode, ok := e.nodeMaps[to.uuid]
	if !ok {
		return fmt.Errorf("to actor not found")
	}

	err := e.DAG.AddEdge(fromNode, toNode)
	if err != nil {
		e.logger.Fatalw("Error adding edge", "error", err)
		return err
	}

	// TODO: add sub and pub to the upstream and downstream actors
	return nil
}

// Ready is the method that generates the DAG, and verifies that the DAG is valid
func (e *Engine[Actor]) Ready() error {
	roots, err := e.getRootActors()
	if err != nil {
		e.logger.Fatalw("Error getting root actor", "error", err)
		return err
	}

	// only one root actor is allowed
	if len(roots) != 1 {
		return fmt.Errorf("only one root actor is allowed")
	}

	e.root = roots[0].Value
	// TODO: active the actors

	return nil
}

//

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

// getRootActor returns the root actors of the DAG, could be multiple
func (e *Engine[Actor]) getRootActors() ([]*pkg.Node[Actor], error) {
	components := e.DAG.StronglyConnectedComponents()
	if len(components) == 0 {
		return nil, fmt.Errorf("No root actor found")
	}

	var roots []*pkg.Node[Actor]
	for _, component := range components {
		for _, node := range component {
			roots = append(roots, node)
		}
	}
	return roots, nil
}
