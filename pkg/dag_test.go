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

	assert.Nil(t, dag.AddEdge(nodeA, nodeB))
	assert.Nil(t, dag.AddEdge(nodeB, nodeC))
	assert.Nil(t, dag.AddEdge(nodeC, nodeD))
	assert.NotNil(t, dag.AddEdge(nodeD, nodeA))

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

// TestDAG leaf nodes
func TestDAG_LeafNodes(t *testing.T) {
	dag := newTestDag(t)
	LeafNodes := dag.LeafNodes()
	assert.Equal(t, 1, len(LeafNodes))
	assert.Equal(t, "F", LeafNodes[0].Value.String())
}

// TestDAG non-leaf nodes
func TestDAG_NonLeafNodes(t *testing.T) {
	dag := newTestDag(t)
	nonLeafNodes := dag.GetNonLeafNodes()
	assert.Equal(t, 6, len(nonLeafNodes))
	assert.Equal(t, "A", nonLeafNodes[0].Value.String())
	assert.Equal(t, "B", nonLeafNodes[1].Value.String())
	assert.Equal(t, "C", nonLeafNodes[2].Value.String())
	assert.Equal(t, "D", nonLeafNodes[3].Value.String())
	assert.Equal(t, "E", nonLeafNodes[4].Value.String())
	assert.Equal(t, "G", nonLeafNodes[5].Value.String())
}

// TestDAG only one node
func TestDAG_OnlyOneNode(t *testing.T) {
	dag := &DAG[Stringer]{}
	nodeA := dag.AddNode(newTestNode("A"))
	nonLeafNodes := dag.GetNonLeafNodes()
	assert.Equal(t, 0, len(nonLeafNodes))
	leafNodes := dag.LeafNodes()
	assert.Equal(t, 1, len(leafNodes))
	assert.Equal(t, "A", leafNodes[0].Value.String())
	assert.Equal(t, "A", nodeA.Value.String())
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
