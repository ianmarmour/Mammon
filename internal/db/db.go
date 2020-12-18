package db

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
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

// Writes the graph to disk as a binary file
func (g *Graph) Write(path string) {
	now := time.Now()
	sec := now.Unix()
	filename := fmt.Sprintf("%s%v.gob", path, sec)

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("Couldn't open file for writing")
	}
	defer f.Close()

	dataEncoder := gob.NewEncoder(f)
	dataEncoder.Encode(g)
}

// Node Represents a node in our graph containing a value
type Node struct {
	ID    int64
	Value interface{}
}

// Load Loads a Graph object from the specified path
func Load(path string) *Graph {
	var data Graph

	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	decoder := gob.NewDecoder(f)

	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	return &data
}
