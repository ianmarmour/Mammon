package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// MediaEntry Represents an item media cache entry.
type MediaEntry struct {
	URL string `json:"url"`
}

// MediaCache Represents a list of entries.
type MediaCache struct {
	Entries map[int64]MediaEntry
}

// Read Reads an existing entry from the cache.
func (c *MediaCache) Read(ID int64) (*MediaEntry, error) {
	if val, ok := c.Entries[ID]; ok {
		return &val, nil
	}

	return nil, errors.New("Entry not found")
}

// Update Either adds a new entry or updates an exisiting Cache entry
func (c *MediaCache) Update(ID int64, entry MediaEntry) *MediaEntry {
	c.Entries[ID] = entry

	return &entry
}

// Delete Removes an entry from the cache by ID
func (c *MediaCache) Delete(ID int64) error {
	if _, ok := c.Entries[ID]; ok {
		delete(c.Entries, ID)
		return nil
	}

	log.Println("Error removing non-existent entry from cache")
	return errors.New("Cannot remove non-existant entry from cache")
}

// Initialize Reads the cache from disk and initializes it
func (c *MediaCache) Initialize() error {
	var initCache map[int64]MediaEntry

	osFile, err := os.Open("item_media.json")
	if err != nil {
		fmt.Println(err)
	}

	osFileBytes, err := ioutil.ReadAll(osFile)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(osFileBytes, &initCache)

	c.Entries = initCache

	defer osFile.Close()

	return nil
}

// Persist Writes the cache to persistant storage
func (c *MediaCache) Persist() error {
	jc, _ := json.Marshal(c.Entries)
	err := ioutil.WriteFile("item_media.json", jc, 0644)

	return err
}
