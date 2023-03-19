package pkg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
	dag := newTestDag(t)
	dag.PrintWithArrows()
}

func TestDAG_StronglyConnectedComponents(t *testing.T) {
	dag := newTestDag(t)
	components := dag.StronglyConnectedComponents()

	for i, component := range components {
		fmt.Printf("Component %d:\n", i+1)
		for _, node := range component {
			fmt.Printf("  %s\n", node.Value.String())
		}
	}
}

func TestDAG_TopologicalSort(t *testing.T) {
	dag := newTestDag(t)
	sortedNodes, err := dag.TopologicalSort()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Topological order:")
	for _, node := range sortedNodes {
		fmt.Printf("  %s\n", node.Value.String())
	}

	assert.Equal(t, 7, len(sortedNodes))
	assert.Equal(t, "G", sortedNodes[0].Value.String())
}

func TestDAG_Neighbors(t *testing.T) {
	dag := newTestDag(t)
	nodeA := dag.Nodes[0]
	neighbors := dag.Neighbors(nodeA)
	assert.Equal(t, 2, len(neighbors))
	assert.Equal(t, "B", neighbors[0].Value.String())
	assert.Equal(t, "C", neighbors[1].Value.String())
}

func newTestDag(t *testing.T) *DAG[Stringer] {

	dag := &DAG[Stringer]{}

	nodeA := dag.AddNode(newTestNode("A"))
	nodeB := dag.AddNode(newTestNode("B"))
	nodeC := dag.AddNode(newTestNode("C"))
	nodeD := dag.AddNode(newTestNode("D"))
	nodeE := dag.AddNode(newTestNode("E"))
	nodeF := dag.AddNode(newTestNode("F"))
	nodeG := dag.AddNode(newTestNode("G"))

	assert.Nil(t, dag.AddEdge(nodeA, nodeB))
	assert.Nil(t, dag.AddEdge(nodeA, nodeC))
	assert.Nil(t, dag.AddEdge(nodeC, nodeD))
	assert.Nil(t, dag.AddEdge(nodeB, nodeE))
	assert.Nil(t, dag.AddEdge(nodeE, nodeF))
	assert.Nil(t, dag.AddEdge(nodeD, nodeF))
	assert.Nil(t, dag.AddEdge(nodeG, nodeB))
	return dag
}
