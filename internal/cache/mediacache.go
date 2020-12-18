package cache

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
)

// MediaEntry Represents an item media cache entry.
type MediaEntry struct {
	URL string
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

// Load Loads a Graph object from the specified path
func Load(path string) *MediaCache {
	var data MediaCache

	f, err := os.Open(path)
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
