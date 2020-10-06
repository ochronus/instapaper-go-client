package instapaper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gomodule/oauth1/oauth"
)

var (
	mux                *http.ServeMux
	client             Client
	server             *httptest.Server
	defaultCredentials = &oauth.Credentials{
		Token:  "Yolo",
		Secret: "Yolo",
	}
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client, _ = NewClient("client_id", "client_secret", "username", "password")
	client.BaseURL = server.URL
	client.Credentials = defaultCredentials
	client.OAuthClient.TokenRequestURI = server.URL + "/oauth/access_token"
}

func teardown() {
	server.Close()
}

func TestWithoutAuthentication(t *testing.T) {
	setup()
	defer teardown()
	client.Credentials = nil
	mux.HandleFunc("/bookmarks/list", func(w http.ResponseWriter, r *http.Request) {
	})
	svc := BookmarkService{
		Client: client,
	}
	_, err := svc.List(DefaultBookmarkListRequestParams)
	if err == nil {
		t.Errorf("expected err NOT to be nil, got %v", err)
	}
}

func TestBogusValidResponse(t *testing.T) {
	setup()
	defer teardown()
	rawResponse := `{"A":"a"}`
	mux.HandleFunc("/bookmarks/list", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodPost; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}
		fmt.Fprint(w, rawResponse)
	})
	svc := BookmarkService{
		Client: client,
	}
	bookmarkList, err := svc.List(DefaultBookmarkListRequestParams)
	if err != nil {
		t.Errorf("expected err to be nil, got %v", err)
	}
	if !reflect.DeepEqual(bookmarkList, &BookmarkListResponse{
		RawResponse: rawResponse,
	}) {
		t.Errorf("Expected the returned bookmark list to only contain the raw response, got %v", bookmarkList)
	}
}

func TestInvalidResponse(t *testing.T) {
	setup()
	defer teardown()
	rawResponse := `qf3434f3g`
	mux.HandleFunc("/bookmarks/list", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, rawResponse)
	})
	svc := BookmarkService{
		Client: client,
	}
	bookmarkList, err := svc.List(DefaultBookmarkListRequestParams)
	if err == nil {
		t.Errorf("expected err NOT to be nil")
	}
	if !reflect.DeepEqual(bookmarkList, &BookmarkListResponse{
		RawResponse: rawResponse,
	}) {
		t.Errorf("Expected the returned bookmark list to to only contain the raw response, got %v", bookmarkList)
	}
}

func TestNot200OKResponse(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/bookmarks/list", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	svc := BookmarkService{
		Client: client,
	}
	bookmarkList, err := svc.List(DefaultBookmarkListRequestParams)
	if err == nil {
		t.Errorf("expected err NOT to be nil")
	}
	if !reflect.DeepEqual(bookmarkList, &BookmarkListResponse{}) {
		t.Errorf("Expected the returned bookmark list to be empty, got %v", bookmarkList)
	}
}

func TestValidResponse(t *testing.T) {
	setup()
	defer teardown()
	rawResponse := `
	{
		"highlights":[
		   {
			  "highlight_id":123456,
			  "text":"That said, I do have some feelings on the matter.",
			  "note":null,
			  "bookmark_id":123456,
			  "time":1601797631,
			  "position":0,
			  "type":"highlight"
		   }
		],
		"bookmarks":[
		   {
			  "hash":"hashety hash",
			  "description":"yo description",
			  "bookmark_id":123456,
			  "private_source":"",
			  "title":"On Call Shouldn\u2019t Suck: A Guide For Managers",
			  "url":"https://charity.wtf/2020/10/03/on-call-shouldnt-suck-a-guide-for-managers/",
			  "progress_timestamp":0,
			  "time":1601750093,
			  "progress":0.0,
			  "starred":"0",
			  "type":"bookmark"
		   }
		],
		"user":{
		   "username":"nope@nope.com",
		   "user_id":12345678,
		   "type":"user",
		   "subscription_is_active":"1"
		}
	 }
	`
	expectedResponse := BookmarkListResponse{
		RawResponse: rawResponse,
		Bookmarks: []Bookmark{
			Bookmark{
				Hash:              "hashety hash",
				Description:       "yo description",
				ID:                123456,
				PrivateSource:     "",
				Title:             "On Call Shouldn\u2019t Suck: A Guide For Managers",
				URL:               "https://charity.wtf/2020/10/03/on-call-shouldnt-suck-a-guide-for-managers/",
				ProgressTimestamp: 0,
				Time:              1601750093,
				Progress:          0.0,
				Starred:           "0",
			},
		},
		Highlights: []Highlight{
			Highlight{
				ID:         123456,
				BookmarkID: 123456,
				Text:       "That said, I do have some feelings on the matter.",
				Note:       "",
				Time:       1601797631,
				Position:   0,
			},
		},
	}

	mux.HandleFunc("/bookmarks/list", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodPost; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}
		fmt.Fprint(w, rawResponse)
	})
	svc := BookmarkService{
		Client: client,
	}
	bookmarkList, err := svc.List(DefaultBookmarkListRequestParams)
	if err != nil {
		t.Errorf("expected err to be nil, got %v", err)
	}
	if !reflect.DeepEqual(*bookmarkList, expectedResponse) {
		t.Errorf("Expected the returned bookmark list to be %v, instead got %v", expectedResponse, bookmarkList)
	}
}
