package api

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ianmarmour/Mammon/pkg/rhttp"
)

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
func getBody(req *http.Request, client *rhttp.RLHTTPClient) ([]byte, error) {
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode > 299 {
		log.Println("Bad status code")
		log.Println(r.StatusCode)
		return nil, err
	}

	var reader io.ReadCloser

	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(r.Body)
		defer reader.Close()

	default:
		reader = r.Body
	}

	body, readErr := ioutil.ReadAll(reader)
	if readErr != nil {
		return nil, readErr
	}

	return body, nil
}
