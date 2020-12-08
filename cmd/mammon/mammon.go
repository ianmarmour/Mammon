package main

import (
	"log"
	"net/http"
	"time"

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

	for _, realm := range res.Realms {
		auctions, err := client.GetAuctions(realm.ID)
		if err != nil {
			log.Println(err)
		}

		log.Println(auctions)
	}
}
