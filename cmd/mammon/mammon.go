package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ianmarmour/Mammon/internal/db"
	"github.com/ianmarmour/Mammon/pkg/blizzard"
	"github.com/ianmarmour/Mammon/pkg/config"
)

func main() {
	config, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	http := &http.Client{Timeout: 10 * time.Second}
	client := &blizzard.Client{nil, *config, http}

	err = client.Authenticate()
	if err != nil {
		log.Fatal(err)
	}

	// Currently disfunctional, need to add logic to fetch or get token from environment.
	res, err := client.GetRealmsIndex()
	if err != nil {
		log.Fatal(err)
	}

	g := db.Graph{}

	for _, realm := range res.Realms {
		rNode := db.Node{}
		rNode.Value = realm

		g.AddNode(rNode)
		log.Println("Added realm to graph")

		auctions, err := client.GetAuctions(realm.ID)
		if err != nil {
			log.Println(err)
		} else {
			for _, auction := range auctions.Auctions {
				log.Println(auction.Item.ID)
				aNode := db.Node{}
				aNode.Value = auction
				g.AddNode(aNode)
				g.AddEdge(rNode, aNode)
			}

			g.Write(config.DBPath)

			os.Exit(0)
		}
	}
}
