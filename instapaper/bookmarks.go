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
	Limit           int
	Skip            []Bookmark
	CustomHaveParam string
	Folder          string
}

// DefaultBookmarkListRequestParams provides sane defaults - no filtering and the maximum limit of 500 bookmarks
var DefaultBookmarkListRequestParams = BookmarkListRequestParams{
	Limit:           500,
	Skip:            nil,
	CustomHaveParam: "",
	Folder:          FolderIDUnread,
}

// BookmarkService defines the interface for all bookmark related API operations
type bookmarkService interface {
	List(BookmarkListRequestParams) ([]Bookmark, error)
	GetText(int) (string, error)
	Star(int) error
	UnStar(int) error
	Archive(int) error
	UnArchive(int) error
	DeletePermanently(int) error
	Move(int, string) error
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
	if p.CustomHaveParam != "" {
		params.Set("have", p.CustomHaveParam)
	} else {
		var haveList []string
		for _, bookmark := range p.Skip {
			haveList = append(haveList, strconv.Itoa(bookmark.ID))
		}
		params.Set("have", strings.Join(haveList, ","))
	}

	res, err := svc.Client.Call("/bookmarks/list", params)
	if err != nil {
		return &BookmarkListResponse{}, err
	} else {
		var bookmarkList BookmarkListResponse
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, &APIError{
				StatusCode:   res.StatusCode,
				Message:      err.Error(),
				ErrorCode:    ErrHTTPError,
				WrappedError: err,
			}
		}
		bodyString := string(bodyBytes)
		bookmarkList.RawResponse = bodyString
		err = json.Unmarshal([]byte(bodyString), &bookmarkList)
		if err != nil {
			return &BookmarkListResponse{
					RawResponse: bodyString,
				}, &APIError{
					StatusCode:   res.StatusCode,
					Message:      err.Error(),
					ErrorCode:    ErrUnmarshalError,
					WrappedError: err,
				}
		}
		return &bookmarkList, nil
	}

}

// GetText returns the specified bookmark's processed text-view HTML, which is always text/html encoded as UTF-8.
func (svc *BookmarkServiceOp) GetText(bookmarkID int) (string, error) {
	params := url.Values{}
	params.Set("bookmark_id", strconv.Itoa(bookmarkID))
	res, err := svc.Client.Call("/bookmarks/get_text", params)
	if err != nil {
		return "", err
	} else {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", &APIError{
				StatusCode:   res.StatusCode,
				Message:      err.Error(),
				ErrorCode:    ErrHTTPError,
				WrappedError: err,
			}
		}
		return string(bodyBytes), nil
	}
}

// Star stars the specified bookmark
func (svc *BookmarkServiceOp) Star(bookmarkID int) error {
	params := url.Values{}
	params.Set("bookmark_id", strconv.Itoa(bookmarkID))
	_, err := svc.Client.Call("/bookmarks/star", params)
	return err
}

// UnStar un-stars the specified bookmark
func (svc *BookmarkServiceOp) UnStar(bookmarkID int) error {
	params := url.Values{}
	params.Set("bookmark_id", strconv.Itoa(bookmarkID))
	_, err := svc.Client.Call("/bookmarks/unstar", params)
	return err
}

// Archive archives the specified bookmark
func (svc *BookmarkServiceOp) Archive(bookmarkID int) error {
	params := url.Values{}
	params.Set("bookmark_id", strconv.Itoa(bookmarkID))
	_, err := svc.Client.Call("/bookmarks/archive", params)
	return err
}

// UnArchive un-archives the specified bookmark
func (svc *BookmarkServiceOp) UnArchive(bookmarkID int) error {
	params := url.Values{}
	params.Set("bookmark_id", strconv.Itoa(bookmarkID))
	_, err := svc.Client.Call("/bookmarks/unarchive", params)
	return err
}

// DeletePermanently PERMANENTLY deletes the specified bookmark
func (svc *BookmarkServiceOp) DeletePermanently(bookmarkID int) error {
	params := url.Values{}
	params.Set("bookmark_id", strconv.Itoa(bookmarkID))
	_, err := svc.Client.Call("/bookmarks/delete", params)
	return err
}

// Move moves the specified bookmark to the specified folder
func (svc *BookmarkServiceOp) Move(bookmarkID int, folderID string) error {
	params := url.Values{}
	params.Set("bookmark_id", strconv.Itoa(bookmarkID))
	params.Set("folder_id", folderID)
	_, err := svc.Client.Call("/bookmarks/move", params)
	return err
}
