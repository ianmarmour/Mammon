package db

import (
	"fmt"
	"sync"
)

// Graph Represents a Mammon graph database storing all auction related information
type Graph struct {
	nodes []*Node
	edges map[Node][]*Node
	lock  sync.RWMutex
}

// AddNode adds a node to the graph safely
func (g *Graph) AddNode(n Node) {
	g.lock.Lock()
	g.nodes = append(g.nodes, &n)
	g.lock.Unlock()
}

// AddEdge Adds an edge between two nodes
func (g *Graph) AddEdge(n1, n2 Node) {
	g.lock.Lock()

	if g.edges == nil {
		g.edges = make(map[Node][]*Node)
	}

	g.edges[n1] = append(g.edges[n1], &n2)
	g.edges[n2] = append(g.edges[n2], &n1)

	g.lock.Unlock()
}

// GetNeighborhood Gets all the neighbors of a specified node
func (g *Graph) GetNeighborhood(n Node) []*Node {
	edges := g.edges[n]

	return edges
}

// AddEdge adds an edge to the graph
func (g *Graph) String() {
	g.lock.RLock()
	s := ""
	for i := 0; i < len(g.nodes); i++ {
		s += g.nodes[i].String() + " -> "
		near := g.edges[*g.nodes[i]]
		for j := 0; j < len(near); j++ {
			s += near[j].String() + " "
		}
		s += "\n"
	}
	fmt.Println(s)
	g.lock.RUnlock()
}

// Node Represents a node in our graph containing a value
type Node struct {
	ID    int64
	Value interface{}
}
