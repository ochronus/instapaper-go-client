package instapaper

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/gomodule/oauth1/oauth"
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
	Bookmarks  []Bookmark
	Highlights []Highlight
}

type BookmarkListRequestParams struct {
	Limit  int
	Skip   []Bookmark
	Folder string
}

var DefaultBookmarkListRequestParams = BookmarkListRequestParams{
	Limit:  500,
	Skip:   nil,
	Folder: FolderIDUnread,
}

type BookmarkService interface {
	List(BookmarkListRequestParams) ([]Bookmark, error)
}

type BookmarkServiceOp struct {
	Client      oauth.Client
	Credentials *oauth.Credentials
}

func (svc *BookmarkServiceOp) List(p BookmarkListRequestParams) (*BookmarkListResponse, error) {
	params := url.Values{}
	params.Set("limit", strconv.Itoa(p.Limit))
	var haveList []string
	for _, bookmark := range p.Skip {
		haveList = append(haveList, strconv.Itoa(bookmark.ID))
	}
	params.Set("have", strings.Join(haveList, ","))
	res, err := svc.Client.Post(nil, svc.Credentials, "https://www.instapaper.com/api/1.1/bookmarks/list", params)
	if err == nil {
		var bookmarkList BookmarkListResponse
		err := json.NewDecoder(res.Body).Decode(&bookmarkList)
		if err != nil {
			return &BookmarkListResponse{}, err
		}
		return &bookmarkList, nil
	}
	return nil, err
}
