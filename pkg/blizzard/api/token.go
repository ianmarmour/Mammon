package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/ianmarmour/Mammon/pkg/config"
)

// Token Blizzard OAuth Token Response
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json"token_type"`
	ExpiresIn   int64  `json"expires_in"`
	Scope       string `json:"scope"`
}

// GetToken Returns a valid OAuth session token based on Blizzard OAuth app credentials
func GetToken(config config.Config, client *http.Client) (*Token, error) {
	url := fmt.Sprintf("https://%s.%s/oauth/token", config.Region.ID, config.AuthEndpoint)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	label, err := w.CreateFormField("grant_type")
	if err != nil {
		return nil, err
	}
	label.Write([]byte("client_credentials"))
	w.Close()

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.SetBasicAuth(config.Credential.ID, config.Credential.Secret)

	resBody, err := getBody(req, client)
	if err != nil {
		return nil, err
	}

	token := Token{}
	err = json.Unmarshal(resBody, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
