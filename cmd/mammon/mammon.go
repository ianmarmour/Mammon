package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ianmarmour/Mammon/internal/cache"
	"github.com/ianmarmour/Mammon/internal/db"
	"github.com/ianmarmour/Mammon/pkg/blizzard"
	"github.com/ianmarmour/Mammon/pkg/blizzard/api"
	"github.com/ianmarmour/Mammon/pkg/config"
	"github.com/ianmarmour/Mammon/pkg/rhttp"
	"golang.org/x/time/rate"
)

func main() {
	config, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	done := make(chan bool)

	for {
		select {
		case <-ticker.C:
			rl := rate.NewLimiter(80, 80) // 90 TPS limit per Blizzard API
			http := rhttp.NewClient(rl)
			blizzclient := &blizzard.Client{nil, *config, http}

			err = blizzclient.Authenticate()
			if err != nil {
				log.Fatal(err)
			}

			g := db.Graph{}
			mc := cache.Initialize(config.CachePath)
			realms, err := blizzclient.GetConnectedRealmsIndex()
			if err != nil {
				log.Fatal(err)
			}

			getGraphData(blizzclient, &g, realms)
			getCacheData(blizzclient, &g, mc)

			mc.Persist(config.CachePath)
			g.Persist(config.DBPath)
		case <-done:
			return
		}
	}
}

func getGraphData(client *blizzard.Client, graph *db.Graph, realms *api.ConnectedRealmsIndex) {
	var wg sync.WaitGroup

	for _, id := range realms.IDs() {
		wg.Add(1)

		go func(api *blizzard.Client, g *db.Graph, id int64, wg *sync.WaitGroup) {
			defer wg.Done()

			realm, err := api.GetConnectedRealm(id)
			if err != nil {
				msg := fmt.Sprintf("Error fetching connnected realm information for realm with ID: %v", id)
				log.Println(msg)
				return
			}

			auctions, err := api.GetAuctions(id)
			if err != nil {
				msg := fmt.Sprintf("Error fetching auctions for connected realm with ID: %v", id)
				log.Println(msg)
				return
			}

			err = g.PopulateRealm(realm, auctions)
			if err != nil {
				msg := fmt.Sprintf("Error population auctions for connected realm with ID: %v", id)
				log.Println(msg)
			}
		}(client, graph, id, &wg)
	}

	wg.Wait()
}

func getCacheData(client *blizzard.Client, graph *db.Graph, mc *cache.MediaCache) {
	var wg sync.WaitGroup

	// We retrieve pointers to all our graph nodes that represent realms for future cache processing of related auctions.
	rNodes, err := graph.GetRealms()
	if err != nil {
		msg := fmt.Sprintf("Error attempting to get all realm nodes in the DB.")
		log.Fatal(msg)
	}

	ids := make(map[int64]bool)

	for _, node := range rNodes {
		neighbors := graph.GetNeighborhood(node)

		for _, aNode := range neighbors {
			auction := aNode.Value.(api.Auction)

			if mc.Exists(auction.Item.ID) != true {
				ids[auction.Item.ID] = true
			}
		}
	}

	for id := range ids {
		wg.Add(1)

		go func(api *blizzard.Client, mc *cache.MediaCache, itemID int64, wg *sync.WaitGroup) {
			defer wg.Done()

			media, err := client.GetItemMedia(itemID)
			if err != nil {
				msg := fmt.Sprintf("Error attempting to fetch Item Media for Item ID: %v", itemID)
				log.Println(msg)
				return
			}

			mc.Update(itemID, *media)
		}(client, mc, id, &wg)
	}

	wg.Wait()
}
