package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ianmarmour/Mammon/internal/cache"
	"github.com/ianmarmour/Mammon/internal/db"
	"github.com/ianmarmour/Mammon/pkg/blizzard"
	"github.com/ianmarmour/Mammon/pkg/config"
	"github.com/ianmarmour/Mammon/pkg/rhttp"
	"golang.org/x/time/rate"
)

func main() {
	config, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	rl := rate.NewLimiter(rate.Every(1*time.Second), 90) // 90 TPS limit per Blizzard API
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

		go func(api *blizzard.Client, g *db.Graph, crID int64, wg *sync.WaitGroup) {
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

			err = db.PopulateRealm(g, cr, auctions)
			if err != nil {
				msg := fmt.Sprintf("Error population auctions for connected realm with ID: %v", cr.ID)
				log.Println(msg)
			}

		}(blizzclient, &g, crID, &wg)
	}

	wg.Wait()

	mc.Persist(config.CachePath)
	g.Persist(config.DBPath)
	os.Exit(0)
}
