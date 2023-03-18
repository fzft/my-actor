package pkg

import (
	"fmt"
	"testing"
)

type TestNode struct {
	name string
}

func newTestNode(name string) *TestNode {
	return &TestNode{name: name}
}

func (n *TestNode) String() string {
	return n.name
}

func TestDAG_AddEdge(t *testing.T) {
	dag := &DAG[Stringer]{}

	nodeA := dag.AddNode(newTestNode("A"))
	nodeB := dag.AddNode(newTestNode("B"))
	nodeC := dag.AddNode(newTestNode("C"))
	nodeD := dag.AddNode(newTestNode("D"))

	if err := dag.AddEdge(nodeA, nodeB); err != nil {
		fmt.Println("Error adding edge:", err)
	}

	if err := dag.AddEdge(nodeB, nodeC); err != nil {
		fmt.Println("Error adding edge:", err)
	}

	if err := dag.AddEdge(nodeC, nodeD); err != nil {
		fmt.Println("Error adding edge:", err)
	}

	//// This edge would create a cycle, so it should produce an error.
	//if err := dag.AddEdge(nodeD, nodeA); err != nil {
	//	fmt.Println("Error adding edge:", err)
	//}

	dag.PrintWithArrows()
}

func TestDAG_Print(t *testing.T) {
	dag := &DAG[Stringer]{}

	nodeA := dag.AddNode(newTestNode("A"))
	nodeB := dag.AddNode(newTestNode("B"))
	nodeC := dag.AddNode(newTestNode("C"))
	nodeD := dag.AddNode(newTestNode("D"))
	nodeE := dag.AddNode(newTestNode("E"))
	nodeF := dag.AddNode(newTestNode("F"))

	if err := dag.AddEdge(nodeA, nodeB); err != nil {
		fmt.Println("Error adding edge:", err)
	}

	if err := dag.AddEdge(nodeA, nodeC); err != nil {
		fmt.Println("Error adding edge:", err)
	}

	if err := dag.AddEdge(nodeC, nodeD); err != nil {
		fmt.Println("Error adding edge:", err)
	}

	if err := dag.AddEdge(nodeB, nodeE); err != nil {
		fmt.Println("Error adding edge:", err)
	}

	if err := dag.AddEdge(nodeE, nodeF); err != nil {
		fmt.Println("Error adding edge:", err)
	}

	if err := dag.AddEdge(nodeD, nodeF); err != nil {
		fmt.Println("Error adding edge:", err)
	}

	dag.PrintWithArrows()
}
