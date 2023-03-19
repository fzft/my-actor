package pkg

import (
	"fmt"
)

type Stringer interface {
	String() string
}

type TarjanData struct {
	index   int
	lowlink int
	onStack bool
}

type Node[T Stringer] struct {
	Value T
}

func (n *Node[Stringer]) String() string {
	return n.Value.String()
}

type Edge[T Stringer] struct {
	From *Node[T]
	To   *Node[T]
}

type DAG[T Stringer] struct {
	Nodes []*Node[T]
	Edges []*Edge[T]
}

func NewDAG[T Stringer]() *DAG[T] {
	return &DAG[T]{}
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

// StronglyConnectedComponents Tarjan's algorithm
func (dag *DAG[Stringer]) StronglyConnectedComponents() [][]*Node[Stringer] {
	index := 0
	stack := make([]*Node[Stringer], 0)
	data := make(map[*Node[Stringer]]*TarjanData)

	var strongConnect func(node *Node[Stringer]) []*Node[Stringer]
	strongConnect = func(node *Node[Stringer]) []*Node[Stringer] {
		data[node] = &TarjanData{
			index:   index,
			lowlink: index,
			onStack: true,
		}
		index++
		stack = append(stack, node)

		for _, neighbor := range dag.Neighbors(node) {
			if _, found := data[neighbor]; !found {
				strongConnect(neighbor)
				data[node].lowlink = min(data[node].lowlink, data[neighbor].lowlink)
			} else if data[neighbor].onStack {
				data[node].lowlink = min(data[node].lowlink, data[neighbor].index)
			}
		}

		components := make([]*Node[Stringer], 0)
		if data[node].lowlink == data[node].index {
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				data[w].onStack = false
				components = append(components, w)
				if w == node {
					break
				}
			}
		}
		return components
	}

	components := make([][]*Node[Stringer], 0)
	for _, node := range dag.Nodes {
		if _, found := data[node]; !found {
			component := strongConnect(node)
			if len(component) > 0 {
				components = append(components, component)
			}
		}
	}

	return components
}

func (dag *DAG[Stringer]) WeaklyConnectedComponents() [][]*Node[Stringer] {
	visited := make(map[*Node[Stringer]]bool)
	components := make([][]*Node[Stringer], 0)

	var dfs func(node *Node[Stringer], component []*Node[Stringer]) []*Node[Stringer]
	dfs = func(node *Node[Stringer], component []*Node[Stringer]) []*Node[Stringer] {
		if visited[node] {
			return component
		}
		visited[node] = true
		component = append(component, node)

		for _, neighbor := range dag.Neighbors(node) {
			component = dfs(neighbor, component)
		}

		return component
	}

	for _, node := range dag.Nodes {
		if !visited[node] {
			component := dfs(node, []*Node[Stringer]{})
			components = append(components, component)
		}
	}

	return components
}

func (dag *DAG[Stringer]) TopologicalSort() ([]*Node[Stringer], error) {
	visited := make(map[*Node[Stringer]]bool)
	stack := make([]*Node[Stringer], 0)

	var visit func(node *Node[Stringer]) error
	visit = func(node *Node[Stringer]) error {
		if visited[node] {
			return nil
		}
		visited[node] = true

		for _, neighbor := range dag.Neighbors(node) {
			if dag.HasPath(neighbor, node) {
				return fmt.Errorf("cycle detected, topological sort is not possible")
			}
			if err := visit(neighbor); err != nil {
				return err
			}
		}

		stack = append(stack, node)
		return nil
	}

	for _, node := range dag.Nodes {
		if err := visit(node); err != nil {
			return nil, err
		}
	}

	// Reverse the stack to get the correct topological order
	for i, j := 0, len(stack)-1; i < j; i, j = i+1, j-1 {
		stack[i], stack[j] = stack[j], stack[i]
	}

	return stack, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
