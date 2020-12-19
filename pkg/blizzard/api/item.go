package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ianmarmour/Mammon/pkg/config"
)

// Asset Represents a media asset of a particular WoW item
type Asset struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	FileDataID int64  `json:"file_data_id"`
}

// ItemMedia Represents the media for a particular item from WoW API
type ItemMedia struct {
	Links  Links   `json:"_links"`
	Assets []Asset `json:"assets"`
	ID     int64   `json:"id"`
}

// GetItemMedia Retrives the media for a particular item
func GetItemMedia(ItemID int64, config *config.Config, token string, client *http.Client) (*ItemMedia, error) {
	url := fmt.Sprintf("https://%s.%s/data/wow/media/item/%d?namespace=static-%s&locale=%s&access_token=%s", config.Region.ID, config.Endpoint, ItemID, config.Region.ID, config.Locale.ID, token)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resBody, err := getBody(req, client)
	if err != nil {
		return nil, err
	}

	itemMedia := ItemMedia{}
	err = json.Unmarshal(resBody, &itemMedia)
	if err != nil {
		log.Println(ItemID)
		log.Println(url)
		return nil, err
	}

	return &itemMedia, nil
}
