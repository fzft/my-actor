package pkg

import (
	"fmt"
)

type Stringer interface {
	String() string
}

type Node[T Stringer] struct {
	Value T
}

type Edge[T Stringer] struct {
	From *Node[T]
	To   *Node[T]
}

type DAG[T Stringer] struct {
	Nodes []*Node[T]
	Edges []*Edge[T]
}

func (dag *DAG[Stringer]) AddNode(value Stringer) *Node[Stringer] {
	node := &Node[Stringer]{Value: value}
	dag.Nodes = append(dag.Nodes, node)
	return node
}

func (dag *DAG[Stringer]) AddEdge(from, to *Node[Stringer]) error {
	if dag.HasPath(to, from) {
		return fmt.Errorf("adding this edge would create a cycle")
	}
	edge := &Edge[Stringer]{From: from, To: to}
	dag.Edges = append(dag.Edges, edge)
	return nil
}

func (dag *DAG[Stringer]) HasPath(from, to *Node[Stringer]) bool {
	visited := make(map[*Node[Stringer]]bool)
	var visit func(current *Node[Stringer]) bool
	visit = func(current *Node[Stringer]) bool {
		if current == to {
			return true
		}
		if visited[current] {
			return false
		}
		visited[current] = true
		for _, neighbor := range dag.Neighbors(current) {
			if visit(neighbor) {
				return true
			}
		}
		return false
	}
	return visit(from)
}

func (dag *DAG[Stringer]) Neighbors(node *Node[Stringer]) []*Node[Stringer] {
	neighbors := make([]*Node[Stringer], 0)
	for _, edge := range dag.Edges {
		if edge.From == node {
			neighbors = append(neighbors, edge.To)
		}
	}
	return neighbors
}

func (dag *DAG[T]) PrintWithArrows() {
	for _, edge := range dag.Edges {
		from := edge.From.Value
		to := edge.To.Value
		fmt.Printf("%v -> %v\n", from, to)
	}
}
