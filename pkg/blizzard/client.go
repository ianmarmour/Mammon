package blizzard

import (
	"net/http"

	"github.com/ianmarmour/Mammon/pkg/blizzard/api"
	"github.com/ianmarmour/Mammon/pkg/config"
)

type Client struct {
	Token      *api.Token
	Config     config.Config
	HTTPClient *http.Client
}

func (c *Client) Authenticate() error {
	res, err := api.GetToken(c.Config, c.HTTPClient)
	if err != nil {
		return err
	}

	c.Token = res

	return nil
}

func (c *Client) GetRealmsIndex() (*api.RealmsIndex, error) {
	res, err := api.GetRealmsIndex(&c.Config, c.Token.AccessToken, c.HTTPClient)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) GetAuctions(realmID int64) (*api.Auctions, error) {
	res, err := api.GetAuctions(realmID, &c.Config, c.Token.AccessToken, c.HTTPClient)
	if err != nil {
		return nil, err
	}

	return res, nil
}
