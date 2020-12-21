package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/ianmarmour/Mammon/pkg/config"
	"github.com/ianmarmour/Mammon/pkg/rhttp"
)

type Status struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Population struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Region struct {
	Key  Link   `json:"key"`
	Name string `json:"name"`
	ID   int64  `json:"ID"`
}

type BType struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type ConnectedRealm struct {
	Links              Links      `json:"_links"`
	ID                 int64      `json:"id"`
	HasQueue           bool       `json:"has_queue"`
	Status             Status     `json:"status"`
	Population         Population `json:"population"`
	Realms             []Realm    `json:"realms"`
	MythicLeaderboards Link       `json:"mythic_leaderboards"`
	Auctions           Link       `json:"auctions"`
}

type Realm struct {
	ID             int64  `json:"id"`
	Region         Region `json:"region"`
	ConnectedRealm Link   `json:"connected_realm"`
	Name           string `json:"name"`
	Category       string `json:"category"`
	Locale         string `json:"locale"`
	TimeZone       string `json:"timezone"`
	Type           BType  `json:"type"`
	IsTournament   bool   `json:"is_tournament"`
	Slug           string `json:"slug"`
}

// RealmIndexEntry Represents a World of Warcraft realm
type RealmIndexEntry struct {
	Key  Link   `json:"key"`
	Name string `json:"name"`
	ID   int64  `json:"id"`
	Slug string `json:"slug"`
}

// RealmsIndex Contains an index of all existing World of Warcraft Realms
type RealmsIndex struct {
	Links  Links             `json:"_links"`
	Realms []RealmIndexEntry `json:"realms"`
}

// IDs returns the IDs of all realms in the realms index
func (r *RealmsIndex) IDs() []int64 {
	var ids []int64

	for _, realm := range r.Realms {
		ids = append(ids, realm.ID)
	}

	return ids
}

// ConnectedRealmsIndex Contains an index of all existing World of Warcraft Connected Realms
type ConnectedRealmsIndex struct {
	Links           Links  `json:"_links"`
	ConnectedRealms []Link `json:"connected_realms"`
}

// IDs returns the IDs of all realms in the realms index
func (r *ConnectedRealmsIndex) IDs() []int64 {
	var ids []int64
	reg, err := regexp.Compile("(?:connected-realm\\/)(.*)(?:\\?)")
	if err != nil {
		log.Fatal("Error constructing regex to parse realmIDs from ConnectedRealm Href")
	}

	for _, realm := range r.ConnectedRealms {
		// Annoying we can't optimize this because golang regex implementation doesn't support negative lookarounds... ree
		match := reg.FindStringSubmatch(realm.Href)
		i64, err := strconv.ParseInt(match[1], 10, 32)
		if err != nil {
			log.Fatal("Error converting connected Realm ID string to int64 exiting...")
		}

		ids = append(ids, i64)
	}

	return ids
}

// GetRealmsIndex Retrives the realms index in a given region with a provided locale
func GetRealmsIndex(config *config.Config, token string, client *rhttp.RLHTTPClient) (*RealmsIndex, error) {
	url := fmt.Sprintf("https://%s.%s/data/wow/realm/index?namespace=dynamic-%s&locale=%s&access_token=%s", config.Region.ID, config.Endpoint, config.Region.ID, config.Locale.ID, token)

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

// GetConnectedRealmsIndex Retrives the connected realms index in a given region with a provided locale
func GetConnectedRealmsIndex(config *config.Config, token string, client *rhttp.RLHTTPClient) (*ConnectedRealmsIndex, error) {
	url := fmt.Sprintf("https://%s.%s/data/wow/connected-realm/index?namespace=dynamic-%s&locale=%s&access_token=%s", config.Region.ID, config.Endpoint, config.Region.ID, config.Locale.ID, token)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resBody, err := getBody(req, client)
	if err != nil {
		return nil, err
	}

	connectedRealmsIndex := ConnectedRealmsIndex{}
	err = json.Unmarshal(resBody, &connectedRealmsIndex)
	if err != nil {
		return nil, err
	}

	return &connectedRealmsIndex, nil
}

// GetConnectedRealm Retrives the connected realm based on an ID
func GetConnectedRealm(id int64, config *config.Config, token string, client *rhttp.RLHTTPClient) (*ConnectedRealm, error) {
	url := fmt.Sprintf("https://%s.%s/data/wow/connected-realm/%v?namespace=dynamic-%s&locale=%s&access_token=%s", config.Region.ID, config.Endpoint, id, config.Region.ID, config.Locale.ID, token)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resBody, err := getBody(req, client)
	if err != nil {
		return nil, err
	}

	connectedRealm := ConnectedRealm{}
	err = json.Unmarshal(resBody, &connectedRealm)
	if err != nil {
		return nil, err
	}

	return &connectedRealm, nil
}
