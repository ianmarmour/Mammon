package cache

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ianmarmour/Mammon/pkg/blizzard"
	"github.com/ianmarmour/Mammon/pkg/blizzard/api"
)

// MediaCache Represents a list of entries.
type MediaCache struct {
	Entries map[int64]api.ItemMedia
	lock    sync.RWMutex
}

// Exists Check if an entry exists in the cache.
func (c *MediaCache) Exists(ID int64) bool {
	c.lock.Lock()

	if _, ok := c.Entries[ID]; ok {
		c.lock.Unlock()
		return true
	}

	return false
}

// Read Reads an existing entry from the cache.
func (c *MediaCache) Read(ID int64) *api.ItemMedia {
	c.lock.Lock()
	val := c.Entries[ID]
	c.lock.Unlock()

	return &val
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

// Initialize Either loads an existing cache or returns a new fresh mediacache
func Initialize(path string) *MediaCache {
	var c *MediaCache

	cOk := exists(path)
	if cOk == false {
		c = &MediaCache{Entries: map[int64]api.ItemMedia{}}
	} else {
		// TODO: This function should probably return errors as well incase some other outstanding issue exists with the cache.
		c = load(path)
	}

	return c
}

// Populate Populates our media cache with any missing entries based on itemIDs from the Blizzard API.
func Populate(mc *MediaCache, client *blizzard.Client, itemIDs map[int64]bool) {
	var wg sync.WaitGroup

	for ID := range itemIDs {
		if mc.Exists(ID) == false {
			wg.Add(1)
			go populateEntry(&wg, mc, client, ID)
		}
	}

	wg.Wait()
}

// PopulateEntry Populates a cache entry
func populateEntry(wg *sync.WaitGroup, mc *MediaCache, client *blizzard.Client, ID int64) {
	defer wg.Done()

	im, err := client.GetItemMedia(ID)
	if err != nil {
		msg := fmt.Sprintf("Error attempting to populate media cache for ID: %v", ID)
		log.Println(msg)
		log.Println(err)
	} else {
		mc.Update(ID, *im)
	}
}

// exists determines if the cache exists or not
func exists(path string) bool {
	cacheFilePath := fmt.Sprintf("%smedia.gob", path)
	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		return false
	}

	return true
}

// load Loads a cache object from the specified path
func load(path string) *MediaCache {
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
