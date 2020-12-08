package api

import (
	"io/ioutil"
	"net/http"
)

var endpoint = "api.blizzard.com"

// Item Represents a generic Item in Blizzards APIs
type Item struct {
	ID int64 `json:"id"`
}

// Link Represents a generic link in Blizzards APIs
type Link struct {
	Href string `json:"href"`
}

// Links Is present on every API response from Blizzard
type Links struct {
	Self Link `json:"self"`
}

// Gets the byte data out of a req body.
func getBody(req *http.Request, client *http.Client) ([]byte, error) {
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if r.StatusCode > 400 {
		return nil, err
	}
	defer r.Body.Close()

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		return nil, readErr
	}

	return body, nil
}