package instapaper

import (
	"github.com/gomodule/oauth1/oauth"
)

const defaultBaseURL = "https://www.instapaper.com/api/1.1"

// Client represents the API client and is used directly by the specific API endpoint client implementations
type Client struct {
	OAuthClient oauth.Client
	Credentials *oauth.Credentials
	BaseURL     string
}

// NewClient configures a new Client and returns it. This is the preferred way to get a new client.
func NewClient(oauthClient oauth.Client, credentials *oauth.Credentials) (Client, error) {
	return Client{
		OAuthClient: oauthClient,
		Credentials: credentials,
		BaseURL:     defaultBaseURL,
	}, nil
}
