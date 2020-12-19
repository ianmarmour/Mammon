package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ianmarmour/Mammon/internal/cache"
	"github.com/ianmarmour/Mammon/internal/db"
	"github.com/ianmarmour/Mammon/pkg/blizzard"
	"github.com/ianmarmour/Mammon/pkg/blizzard/api"
	"github.com/ianmarmour/Mammon/pkg/config"
	"go.uber.org/ratelimit"
)

func main() {
	config, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	http := &http.Client{Timeout: 10 * time.Second}
	client := &blizzard.Client{nil, *config, http}

	// Perform OAuth Authentication to Blizzards API and cache response token. Later we should have logic here that can refresh our token on expiration.
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
	mc := mcInit(config.CachePath)

	for _, realm := range res.Realms {
		rNode := db.Node{}
		rNode.Value = realm

		g.AddNode(rNode)
		log.Println("Added realm to graph")

		itemIDs := make(map[int64]bool)

		auctions, err := client.GetAuctions(realm.ID)
		if err != nil {
			log.Println(err)
		} else {
			for _, auction := range auctions.Auctions {
				aNode := db.Node{}
				aNode.Value = auction
				g.AddNode(aNode)
				g.AddEdge(rNode, aNode)

				// Setup for future media cache population
				itemIDs[auction.Item.ID] = true
			}

			err = mcPopulate(mc, client, itemIDs)
			if err != nil {
				log.Println("Fatal error populating media cache")
				os.Exit(1)
			}

			os.Exit(0)
		}
	}

	g.Persist(config.DBPath)
	os.Exit(0)
}

// mcPopuplate Populates our media cache with any missing entries based on itemIDs
func mcPopulate(mc *cache.MediaCache, client *blizzard.Client, itemIDs map[int64]bool) error {
	var wg sync.WaitGroup
	var itemsMedia []*api.ItemMedia
	ch := make(chan *api.ItemMedia)
	rl := ratelimit.New(90)

	for ID := range itemIDs {
		log.Println(ID)
		wg.Add(1)
		rl.Take()
		go getItemMedia(ch, &wg, mc, client, ID)
	}

	wg.Wait()
	close(ch)

	log.Println(itemsMedia)

	return nil
}

// getItemMedia Checks cache for item media and if it doesn't exist gets the media from Blizzard
func getItemMedia(ch chan *api.ItemMedia, wg *sync.WaitGroup, mc *cache.MediaCache, client *blizzard.Client, ID int64) {
	defer wg.Done()

	var im *api.ItemMedia

	im, err := mc.Read(ID)
	if err != nil {
		im, err = client.GetItemMedia(ID)
		if err != nil {
			log.Println(err)
		} else {
			mc.Update(ID, *im)

			ch <- im
		}
	}

	ch <- im
}

// Either loads or sets up the mediaCache
func mcInit(path string) *cache.MediaCache {
	var c *cache.MediaCache

	cOk := cache.Exists(path)
	if cOk == false {
		c = &cache.MediaCache{Entries: map[int64]api.ItemMedia{}}
	} else {
		// TODO: This function should probably return errors as well incase some other outstanding issue exists with the cache.
		c = cache.Load(path)
	}

	return c
}
