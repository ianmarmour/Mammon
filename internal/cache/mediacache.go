package cache

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ianmarmour/Mammon/pkg/blizzard/api"
)

// MediaCache Represents a list of entries.
type MediaCache struct {
	Entries map[int64]api.ItemMedia
	lock    sync.RWMutex
}

// Read Reads an existing entry from the cache.
func (c *MediaCache) Read(ID int64) (*api.ItemMedia, error) {
	c.lock.Lock()

	if val, ok := c.Entries[ID]; ok {
		c.lock.Unlock()

		return &val, nil
	}

	c.lock.Unlock()

	return nil, errors.New("Entry not found")
}

// Update Either adds a new entry or updates an exisiting Cache entry
func (c *MediaCache) Update(ID int64, entry api.ItemMedia) *api.ItemMedia {
	c.lock.Lock()

	c.Entries[ID] = entry

	c.lock.Unlock()

	return &entry
}

// Delete Removes an entry from the cache by ID
func (c *MediaCache) Delete(ID int64) error {
	c.lock.Lock()

	if _, ok := c.Entries[ID]; ok {
		delete(c.Entries, ID)

		c.lock.Unlock()

		return nil
	}

	c.lock.Unlock()

	log.Println("Error removing non-existent entry from cache")
	return errors.New("Cannot remove non-existant entry from cache")
}

// Persist the graph to disk as a binary file
func (c *MediaCache) Persist(path string) {
	filename := fmt.Sprintf("%smedia.gob", path)

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("Couldn't open file for writing")
	}
	defer f.Close()

	dataEncoder := gob.NewEncoder(f)
	dataEncoder.Encode(c)
}

// Exists determines if the cache exists or not
func Exists(path string) bool {
	cacheFilePath := fmt.Sprintf("%smedia.gob", path)
	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		return false
	}

	return true
}

// Load Loads a Graph object from the specified path
func Load(path string) *MediaCache {
	cacheFilePath := fmt.Sprintf("%smedia.gob", path)

	var data MediaCache

	f, err := os.Open(cacheFilePath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	decoder := gob.NewDecoder(f)

	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	return &data
}
