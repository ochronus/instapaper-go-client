package instapaper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

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

// Call makes a call to the Instapaper API on the specific path with the given call parameters. It handles errors converting them to an APIError instance
func (svc *Client) Call(path string, params url.Values) (*http.Response, error) {
	if svc.Credentials == nil {
		return nil, &APIError{
			Message:   "Please call Authenticate() first",
			ErrorCode: ErrNotAuthenticated,
		}
	}
	res, err := svc.OAuthClient.Post(nil, svc.Credentials, svc.BaseURL+path, params)
	if err == nil && res.StatusCode == 200 {
		return res, nil
	}
	var apiError []APIError
	bodyBytes, err := ioutil.ReadAll(res.Body)
	// there was a "low level" http error
	if err != nil {
		return nil, &APIError{
			StatusCode:   res.StatusCode,
			Message:      err.Error(),
			ErrorCode:    ErrHTTPError,
			WrappedError: err,
		}
	}
	err = json.Unmarshal(bodyBytes, &apiError)
	if err != nil {
		return nil, &APIError{
			StatusCode:   res.StatusCode,
			Message:      err.Error(),
			ErrorCode:    ErrUnmarshalError,
			WrappedError: err,
		}
	}
	apiError[0].StatusCode = res.StatusCode
	return nil, &apiError[0]
}
