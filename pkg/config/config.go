package config

import (
	"errors"
	"fmt"
	"log"
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

func GetDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	path := fmt.Sprintf("%s/.mammon/db/", home)

	return path
}

func GetCachePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	path := fmt.Sprintf("%s/.mammon/cache/", home)

	return path
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

	c.Endpoint = "api.blizzard.com"
	c.AuthEndpoint = "battle.net"
	c.Region = region
	c.Locale = locale
	c.Credential = credential
	c.CachePath = GetCachePath()
	c.DBPath = GetDBPath()

	return &c, nil
}
