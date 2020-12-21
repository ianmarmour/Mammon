package db

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ianmarmour/Mammon/pkg/blizzard/api"
)

// Graph Represents a Mammon graph database storing all auction related information
type Graph struct {
	Nodes []*Node
	Edges map[Node][]*Node
	lock  sync.RWMutex
}

// AddNode adds a node to the graph safely
func (g *Graph) AddNode(n *Node) {
	g.lock.Lock()
	g.Nodes = append(g.Nodes, n)
	g.lock.Unlock()
}

// AddEdge Adds an edge between two nodes
func (g *Graph) AddEdge(n1, n2 *Node) {
	g.lock.Lock()

	if g.Edges == nil {
		g.Edges = make(map[Node][]*Node)
	}

	g.Edges[*n1] = append(g.Edges[*n1], n2)

	g.lock.Unlock()
}

// GetNeighborhood Gets all the neighbors of a specified node
func (g *Graph) GetNeighborhood(n *Node) []*Node {
	g.lock.Lock()
	edges := g.Edges[*n]
	g.lock.Unlock()

	return edges
}

// Persist the graph to disk as a binary file
func (g *Graph) Persist(path string) {
	now := time.Now()
	sec := now.Unix()
	filename := fmt.Sprintf("%s%v.gob", path, sec)

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("Couldn't open file for writing")
	}
	defer f.Close()
	// Register our unknown types here.
	gob.Register(api.ConnectedRealm{})
	gob.Register(api.Auction{})

	dataEncoder := gob.NewEncoder(f)
	err = dataEncoder.Encode(g)
	if err != nil {
		log.Println(err)
	}
}

// Node Represents a node in our graph containing a value
type Node struct {
	ID    int64
	Type  string
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

// GetRealms Returns all realm Nodes in the graph.
func (g *Graph) GetRealms() ([]*Node, error) {
	var realms []*Node

	for node := range g.Edges {
		if node.Type == "realm" {
			realms = append(realms, &node)
		}
	}

	return realms, nil
}

// PopulateRealm Populates a realm with all it's auctions in the DB
func (g *Graph) PopulateRealm(cr *api.ConnectedRealm, auctions *api.Auctions) error {
	rNode := Node{Type: "realm"}
	rNode.ID = cr.ID
	rNode.Value = cr
	g.AddNode(&rNode)

	s := fmt.Sprintf("Adding entry for Connected Realm ID: %v to DB", cr.ID)
	log.Println(s)

	for _, auction := range auctions.Auctions {
		aNode := Node{Type: "auction"}
		aNode.ID = auction.ID
		aNode.Value = auction
		g.AddNode(&aNode)
		g.AddEdge(&rNode, &aNode)
	}

	return nil
}
