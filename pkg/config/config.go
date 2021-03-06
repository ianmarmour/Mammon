package config

import (
	"errors"
	"os"
)

// Config Represents a configuration object for Mammon containing all required information to run
type Config struct {
	Endpoint     string
	AuthEndpoint string
	Region       *Region
	Locale       *Locale
	Credential   *Credential
	DBPath       string
	CachePath    string
}

// Region Represents the region for the Blizzard API Endpoint
type Region struct {
	ID string
}

// Locale Represents the proper locale for the Blizzard API
type Locale struct {
	ID string
}

// Cache Represents Cache configuration
type Cache struct {
	Path string
}

// DB Represents DB configuration
type DB struct {
	Path string
}

// Credential Represents the OAuth Credential Pair used to generate the OAuth token for the Blizzard API
type Credential struct {
	ID     string
	Secret string
}

// GetLocale Reads BLIZZARD API Client Locale from the environment
func GetLocale() (*Locale, error) {
	l := Locale{}

	id, idOk := os.LookupEnv("BLIZZARD_API_CLIENT_LOCALE")
	if idOk == false {
		return nil, errors.New("BLIZZARD_API_CLIENT_LOCALE Environmental Variable is Unset")
	}
	l.ID = id

	return &l, nil
}

// GetRegion Reads BLIZZARD API Client Region from the environment
func GetRegion() (*Region, error) {
	r := Region{}

	id, idOk := os.LookupEnv("BLIZZARD_API_CLIENT_REGION")
	if idOk == false {
		return nil, errors.New("BLIZZARD_API_CLIENT_REGION Environmental Variable is Unset")
	}
	r.ID = id

	return &r, nil
}

// GetCredential Reads BLIZZARD OAuth App API Credentials from the environment
func GetCredential() (*Credential, error) {
	c := Credential{}

	id, idOk := os.LookupEnv("BLIZZARD_API_CLIENT_ID")
	if idOk == false {
		return nil, errors.New("BLIZZARD_API_CLIENT_ID Environmental Variable is Unset")
	}
	c.ID = id

	secret, secretOk := os.LookupEnv("BLIZZARD_API_CLIENT_SECRET")
	if secretOk == false {
		return nil, errors.New("BLIZZARD_API_CLIENT_SECRET Environmental Variable is Unset")
	}
	c.Secret = secret

	return &c, nil
}

// GetDB Reads out DB related configuration from the environment
func GetDB() (*DB, error) {
	db := DB{}

	p, pOk := os.LookupEnv("MAMMON_DB_PATH")
	if pOk == false {
		return nil, errors.New("MAMMON_DB_PATH Environmental Variable is Unset")
	}
	db.Path = p

	return &db, nil
}

// GetCache Reads out Cache related configuration from the environment
func GetCache() (*Cache, error) {
	c := Cache{}

	p, pOk := os.LookupEnv("MAMMON_CACHE_PATH")
	if pOk == false {
		return nil, errors.New("MAMMON_CACHE_PATH Environmental Variable is Unset")
	}
	c.Path = p

	return &c, nil
}

func Get() (*Config, error) {
	c := Config{}

	credential, err := GetCredential()
	if err != nil {
		return nil, err
	}

	locale, err := GetLocale()
	if err != nil {
		return nil, err
	}

	region, err := GetRegion()
	if err != nil {
		return nil, err
	}

	cache, err := GetCache()
	if err != nil {
		return nil, err
	}

	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	c.Endpoint = "api.blizzard.com"
	c.AuthEndpoint = "battle.net"
	c.Region = region
	c.Locale = locale
	c.Credential = credential
	c.CachePath = cache.Path
	c.DBPath = db.Path

	return &c, nil
}
