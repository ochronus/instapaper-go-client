package instapaper

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/gomodule/oauth1/oauth"
)

func TestNewClientParams(t *testing.T) {
	newClient, err := NewClient("client_id", "client_secret", "username", "password")
	if err != nil {
		t.Errorf("expected err to be nil, got %v", err)
	}
	if newClient.BaseURL != defaultBaseURL {
		t.Errorf("expected the BaseURL to be equal to defaultBaseURL, %v vs. %v", newClient.BaseURL, defaultBaseURL)
	}
	if !reflect.DeepEqual(newClient.OAuthClient.Credentials, oauth.Credentials{
		Token:  "client_id",
		Secret: "client_secret",
	}) {
		t.Errorf("expected the OAuth client's credentials to be set based on the params, got %v", newClient.OAuthClient.Credentials)
	}
	if newClient.OAuthClient.SignatureMethod != oauth.HMACSHA1 {
		t.Errorf("expected the OAuth client's signature method to be HMAC-SHA1, got %v", newClient.OAuthClient.SignatureMethod)
	}
	if newClient.OAuthClient.TokenRequestURI != defaultBaseURL+"/oauth/access_token" {
		t.Errorf("expected the OAuth client's TokenRequestURI to be %v, got %v", defaultBaseURL+"/oauth/access_token", newClient.OAuthClient.TokenRequestURI)
	}

	if newClient.Username != "username" {
		t.Errorf("expected the username to be 'username', got %v", newClient.Username)
	}
	if newClient.Password != "password" {
		t.Errorf("expected the password to be 'password', got %v", newClient.Password)
	}

}

func TestAuthSuccess(t *testing.T) {
	setup()
	defer teardown()
	creds := oauth.Credentials{
		Token:  "token",
		Secret: "secret",
	}
	mux.HandleFunc("/oauth/access_token", func(w http.ResponseWriter, r *http.Request) {
		params := url.Values{}
		params.Add("oauth_token", creds.Token)
		params.Add("oauth_token_secret", creds.Secret)
		fmt.Fprint(w, params.Encode())
	})
	err := client.Authenticate()
	if err != nil {
		t.Errorf("expected Authenticate() to succed, got %v %v", err, client.OAuthClient.TokenRequestURI)
	}
	if !reflect.DeepEqual(client.Credentials, &creds) {
		t.Errorf("expected the OAuth client's credentials to be %v, got %v", creds, client.Credentials)
	}
}

func TestAuthFailure(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/oauth/access_token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	err := client.Authenticate()
	if err == nil {
		t.Errorf("expected Authenticate() to fail, but it succeded")
	}
	if !reflect.DeepEqual(client.Credentials, defaultCredentials) {
		t.Errorf("expected the OAuth client's credentials to be %v, got %v", defaultCredentials, client.Credentials)
	}
}

func TestErrorHandling(t *testing.T) {
	setup()
	defer teardown()
	rawErrorResponse := `
	[
		{
		   "message":"test error message",
		   "error_code":1337
		}
	 ]
	`
	expectedError := &APIError{
		ErrorCode:    1337,
		Message:      "test error message",
		StatusCode:   http.StatusInternalServerError,
		WrappedError: nil,
	}
	mux.HandleFunc("/errortest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, rawErrorResponse)
	})
	_, err := client.Call("/errortest", nil)
	if !reflect.DeepEqual(err, expectedError) {
		t.Errorf("expected the error to be %v, got %v", expectedError, err)
	}
}
