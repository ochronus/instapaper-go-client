package instapaper

import (
	"github.com/gomodule/oauth1/oauth"
)

const defaultBaseURL = "https://www.instapaper.com/api/1.1"

// Client represents the API client and is used directly by the specific API endpoint client implementations
type Client struct {
	OAuthClient oauth.Client
	Username    string
	Password    string
	Credentials *oauth.Credentials
	BaseURL     string
}

// ClientIf represents the interface an Instapaper API client needs to implement
type ClientIf interface {
	Authenticate() error
}

// NewClient configures a new Client and returns it. This is the preferred way to get a new client.
func NewClient(consumerID string, consumerSecret string, username string, password string) (Client, error) {
	return Client{
		OAuthClient: oauth.Client{
			SignatureMethod: oauth.HMACSHA1,
			Credentials: oauth.Credentials{
				Token:  consumerID,
				Secret: consumerSecret,
			},
			TokenRequestURI: defaultBaseURL + "/oauth/access_token",
		},
		Username: username,
		Password: password,
		BaseURL:  defaultBaseURL,
	}, nil
}

// Authenticate uses the client ID and secret plus the username/password to get oAuth tokens with which it can make authenticated calls in the future
func (svc *Client) Authenticate() error {
	credentials, _, err := svc.OAuthClient.RequestTokenXAuth(nil, nil, svc.Username, svc.Password)
	if err != nil {
		return err
	}
	svc.Credentials = credentials
	return nil
}
