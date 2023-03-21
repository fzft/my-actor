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
	mailbox Mailbox

	// nodeMaps is the map of actors
	nodeMaps map[string]*pkg.Node[T]

	// pidMaps is the map of pid
	pidMaps map[string]*Pid

	// root is the root actor
	root *Pid

	// sinkPool is store the result of the leaf actor
	sinkPool *SinkPool

	isReady bool
}

func NewEngine() *Engine[Actor] {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.TimeKey = ""
	logger, _ := loggerConfig.Build()
	sugarLogger := logger.Sugar()
	return &Engine[Actor]{
		logger:   sugarLogger,
		DAG:      pkg.NewDAG[Actor](),
		mailbox:  NewDefaultMailbox(sugarLogger),
		sinkPool: NewSinkPool(),
		nodeMaps: make(map[string]*pkg.Node[Actor]),
		pidMaps:  make(map[string]*Pid),
	}
}

// Spawn spawns a new actor
func (e *Engine[Actor]) Spawn(actor Actor) (*Pid, error) {
	// check if the actor is already spawned
	if _, ok := e.pidMaps[actor.String()]; ok {
		return nil, fmt.Errorf("actor already spawned")
	}

	pid := NewPid(e.logger, actor)
	node := e.AddNode(actor)
	e.nodeMaps[pid.uuid] = node
	e.pidMaps[actor.String()] = pid
	return pid, nil
}

func (e *Engine[Actor]) AddEdge(from, to *Pid) error {
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

	rootActor := roots[0].Value
	e.root = e.pidMaps[rootActor.String()]

	// setup mailbox to root actor
	e.root.context.setMailbox(e.mailbox)

	// setup every non-leaf actor conn to its child actor
	for _, node := range e.getNonLeafActors() {
		nodePid := e.pidMaps[node.Value.String()]
		childActorNodes := e.getChildActors(node)
		for _, childNode := range childActorNodes {
			childPid := e.pidMaps[childNode.Value.String()]
			nodePid.context.addChild(childPid)
		}
	}

	for _, pid := range e.pidMaps {
		go pid.run()
	}

	// run the tick message
	go e.sinkTickMsg()

	e.isReady = true
	return nil
}

// Send sends a message to the DAG
func (e *Engine[Actor]) Send(msg any) error {
	if !e.isReady {
		return fmt.Errorf("engine is not ready")
	}
	return e.mailbox.Source(msg)
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

// getNonLeafActors returns the non-leaf actors of the DAG, could be multiple
func (e *Engine[Actor]) getNonLeafActors() []*pkg.Node[Actor] {
	return e.DAG.GetNonLeafNodes()
}

// getChildActors returns the child actors of the DAG, could be multiple
func (e *Engine[Actor]) getChildActors(node *pkg.Node[Actor]) []*pkg.Node[Actor] {
	return e.DAG.Neighbors(node)
}

// sinkTickMsg is the message that is sent to the sink actor
func (e *Engine[Actor]) sinkTickMsg() {
	for _, pid := range e.pidMaps {
		go func(pid *Pid) {
			for {
				select {
				case in := <-pid.TickInMsgCh:
					e.sinkPool.PutInMsg(in.uid, in)
				case out := <-pid.TickOutMsgCh:
					e.sinkPool.PutOutMsg(out.uid, out)
				}
			}
		}(pid)
	}
}
