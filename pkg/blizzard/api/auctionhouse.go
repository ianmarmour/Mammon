package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ianmarmour/Mammon/pkg/config"
	"github.com/ianmarmour/Mammon/pkg/rhttp"
)

// Auction represents a single auction returned from the Blizzard auctions API.
type Auction struct {
	ID        int64  `json:"id"`
	Item      Item   `json:"item"`
	Quantity  int64  `json:"quantity"`
	UnitPrice int64  `json:"unit_price"`
	TimeLeft  string `json:"time_left"`
}

// Auctions a list of auctions from a particular realm
type Auctions struct {
	Links          Links     `json:"_links"`
	ConnectedRealm Link      `json:"connected_realm"`
	Auctions       []Auction `json:"auctions"`
}

// ItemIDs returns the Item IDs of all auctions
func (a *Auctions) ItemIDs() []int64 {
	var ids []int64

	for _, auction := range a.Auctions {
		ids = append(ids, auction.ID)
	}

	return ids
}

// GetAuctions Retrives all the active auctions in a particular realm by ID
func GetAuctions(realmID int64, config *config.Config, token string, client *rhttp.RLHTTPClient) (*Auctions, error) {
	url := fmt.Sprintf("https://%s.%s/data/wow/connected-realm/%d/auctions?namespace=dynamic-%s&locale=%s&access_token=%s", config.Region.ID, config.Endpoint, realmID, config.Region.ID, config.Locale.ID, token)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Allow blizzard API to use gzip encoding to speed things up AH data is large.
	req.Header.Add("Accept-Encoding", "gzip")

	resBody, err := getBody(req, client)
	if err != nil {
		return nil, err
	}

	auctions := Auctions{}
	err = json.Unmarshal(resBody, &auctions)
	if err != nil {
		return nil, err
	}

	return &auctions, nil
}
