package instapaper

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
)

// Bookmark represents a single bookmark entry
type Bookmark struct {
	Hash              string
	Description       string
	ID                int    `json:"bookmark_id"`
	PrivateSource     string `json:"private_source"`
	Title             string
	URL               string
	ProgressTimestamp int `json:"progress_timestamp"`
	Time              int
	Progress          float32
	Starred           string
}

// BookmarkListResponse represents the useful part of the API response for the bookmark list endpoint
type BookmarkListResponse struct {
	Bookmarks   []Bookmark
	Highlights  []Highlight
	RawResponse string
}

// BookmarkListRequestParams defines filtering and limiting options for the List endpoint.
// see DefaultBookmarkListRequestParams for a set of sane defaults
type BookmarkListRequestParams struct {
	Limit  int
	Skip   []Bookmark
	Folder string
}

// DefaultBookmarkListRequestParams provides sane defaults - no filtering and the maximum limit of 500 bookmarks
var DefaultBookmarkListRequestParams = BookmarkListRequestParams{
	Limit:  500,
	Skip:   nil,
	Folder: FolderIDUnread,
}

// BookmarkService defines the interface for all bookmark related API operations
type BookmarkService interface {
	List(BookmarkListRequestParams) ([]Bookmark, error)
}

// BookmarkServiceOp is the implementation of the bookmark related parts of the API client, conforming to the BookmarkService interface
type BookmarkServiceOp struct {
	Client Client
}

// List returns the list of bookmarks. By default it returns (maximum) 500 of the unread bookmarks
// see BookmarkListRequestParams for filtering options
func (svc *BookmarkServiceOp) List(p BookmarkListRequestParams) (*BookmarkListResponse, error) {
	params := url.Values{}
	params.Set("limit", strconv.Itoa(p.Limit))
	var haveList []string
	for _, bookmark := range p.Skip {
		haveList = append(haveList, strconv.Itoa(bookmark.ID))
	}
	params.Set("have", strings.Join(haveList, ","))
	url := svc.Client.BaseURL + "/bookmarks/list"
	res, err := svc.Client.OAuthClient.Post(nil, svc.Client.Credentials, url, params)
	if err == nil {
		var bookmarkList BookmarkListResponse
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return &BookmarkListResponse{}, err
		}
		bodyString := string(bodyBytes)
		bookmarkList.RawResponse = bodyString
		err = json.Unmarshal([]byte(bodyString), &bookmarkList)
		if err != nil {
			return &BookmarkListResponse{
				RawResponse: bodyString,
			}, err
		}
		return &bookmarkList, nil
	}
	return nil, err
}
