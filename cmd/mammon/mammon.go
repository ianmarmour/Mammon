package main

import (
	"fmt"
	"log"
	"os"
	"sync"

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

	rl := rate.NewLimiter(80, 80) // 90 TPS limit per Blizzard API
	http := rhttp.NewClient(rl)
	blizzclient := &blizzard.Client{nil, *config, http}

	err = blizzclient.Authenticate()
	if err != nil {
		log.Fatal(err)
	}

	crs, err := blizzclient.GetConnectedRealmsIndex()
	if err != nil {
		log.Fatal(err)
	}

	g := db.Graph{}
	mc := cache.Initialize(config.CachePath)

	var wg sync.WaitGroup

	for _, crID := range crs.IDs() {
		wg.Add(1)

		go func(api *blizzard.Client, g *db.Graph, mc *cache.MediaCache, crID int64, wg *sync.WaitGroup) {
			defer wg.Done()

			cr, err := api.GetConnectedRealm(crID)
			if err != nil {
				msg := fmt.Sprintf("Error fetching connnected realm information for realm with ID: %v", crID)
				log.Println(msg)
				return
			}

			auctions, err := api.GetAuctions(cr.ID)
			if err != nil {
				msg := fmt.Sprintf("Error fetching auctions for connected realm with ID: %v", cr.ID)
				log.Println(msg)
				return
			}

			err = g.PopulateRealm(cr, auctions)
			if err != nil {
				msg := fmt.Sprintf("Error population auctions for connected realm with ID: %v", cr.ID)
				log.Println(msg)
			}
		}(blizzclient, &g, mc, crID, &wg)
	}

	wg.Wait()

	var wg2 sync.WaitGroup

	// We retrieve pointers to all our graph nodes that represent realms for future cache processing of related auctions.
	rNodes, err := g.GetRealms()
	if err != nil {
		msg := fmt.Sprintf("Error attempting to get all realm nodes in the DB.")
		log.Fatal(msg)
	}

	idsToFetch := make(map[int64]bool)

	for _, node := range rNodes {
		neighbors := g.GetNeighborhood(node)

		for _, aNode := range neighbors {
			auction := aNode.Value.(api.Auction)

			if mc.Exists(auction.Item.ID) != true {
				idsToFetch[auction.Item.ID] = true
			}
		}
	}

	for id := range idsToFetch {
		wg2.Add(1)

		go func(api *blizzard.Client, mc *cache.MediaCache, itemID int64, wg *sync.WaitGroup) {
			defer wg2.Done()

			media, err := blizzclient.GetItemMedia(itemID)
			if err != nil {
				msg := fmt.Sprintf("Error attempting to fetch Item Media for Item ID: %v", itemID)
				log.Println(msg)
				return
			}

			mc.Update(itemID, *media)
		}(blizzclient, mc, id, &wg)
	}

	wg2.Wait()

	mc.Persist(config.CachePath)
	g.Persist(config.DBPath)
	os.Exit(0)
}
