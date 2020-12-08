package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Realm Represents a World of Warcraft realm
type Realm struct {
	Key  Link   `json:"key"`
	Name string `json:"name"`
	ID   int64  `json:"id"`
	Slug string `json:"slug"`
}

// RealmsIndex Contains an index of all existing World of Warcraft Realms
type RealmsIndex struct {
	Links  Links   `json:"_links"`
	Realms []Realm `json:"realms"`
}

// GetRealmsIndex Retrives the realms index in a given region with a provided locale
func GetRealmsIndex(token string, region string, locale string, client *http.Client) (*RealmsIndex, error) {
	url := fmt.Sprintf("https://%s.%s/data/wow/realm/index?namespace=dynamic-%s&locale=%s&access_token=%s", region, endpoint, region, locale, token)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resBody, err := getBody(req, client)
	if err != nil {
		return nil, err
	}

	realmsIndex := RealmsIndex{}
	err = json.Unmarshal(resBody, &realmsIndex)
	if err != nil {
		return nil, err
	}

	return &realmsIndex, nil
}
