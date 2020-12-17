package blizzard

import (
	"net/http"

	"github.com/ianmarmour/Mammon/pkg/blizzard/api"
	"github.com/ianmarmour/Mammon/pkg/config"
)

// Client the Blizzard API Client
type Client struct {
	Token      *api.Token
	Config     config.Config
	HTTPClient *http.Client
}

// Authenticate Sets up OAuth communication for the Blizzard API
func (c *Client) Authenticate() error {
	res, err := api.GetToken(c.Config, c.HTTPClient)
	if err != nil {
		return err
	}

	c.Token = res

	return nil
}

// GetRealmsIndex Returns the realm index IDs for all wow realms
func (c *Client) GetRealmsIndex() (*api.RealmsIndex, error) {
	res, err := api.GetRealmsIndex(&c.Config, c.Token.AccessToken, c.HTTPClient)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetAuctions Retrieves all the auctions in a particular realm by ID
func (c *Client) GetAuctions(realmID int64) (*api.Auctions, error) {
	res, err := api.GetAuctions(realmID, &c.Config, c.Token.AccessToken, c.HTTPClient)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetItemMedia Retrives the media metadata for a particular item ID
func (c *Client) GetItemMedia(itemID int64) (*api.ItemMedia, error) {
	res, err := api.GetItemMedia(itemID, &c.Config, c.Token.AccessToken, c.HTTPClient)
	if err != nil {
		return nil, err
	}

	return res, nil
}
