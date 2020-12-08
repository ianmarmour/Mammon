package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ianmarmour/Mammon/pkg/blizzard/api"
)

func main() {
	client := &http.Client{Timeout: 10 * time.Second}
	// Currently disfunctional, need to add logic to fetch or get token from environment.
	res, err := api.GetRealmsIndex("", "us", "en_US", client)
	if err != nil {
		log.Println(err)
	}

	log.Println(res)
}
