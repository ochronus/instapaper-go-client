package instapaper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Bookmark represents a single bookmark entry
type Bookmark struct {
	Hash              string
	Description       string
	ID                int    `json:"bookmark_id"`
	PrivateSource     string `json:"private_source"`
	Title             string
	URL               string
	ProgressTimestamp int64 `json:"progress_timestamp"`
	Time              float32
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
	SkipHighlights  []Highlight
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

// BookmarkAddRequestParams represents all the parameters you can pass when adding a new bookmark
// Either URL or Content or a (Content, PrivateSourceName) pair is mandatory.
type BookmarkAddRequestParams struct {
	URL               string
	Title             string
	Description       string
	Folder            string
	ResolveFinalURL   bool
	Content           string
	PrivateSourceName string
}

// BookmarkService is the implementation of the bookmark related parts of the API client, conforming to the BookmarkService interface
type BookmarkService struct {
	Client Client
}

// List returns the list of bookmarks. By default it returns (maximum) 500 of the unread bookmarks
// see BookmarkListRequestParams for filtering options
func (svc *BookmarkService) List(p BookmarkListRequestParams) (*BookmarkListResponse, error) {
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

	var highlightList []string
	for _, highlight := range p.SkipHighlights {
		highlightList = append(highlightList, strconv.Itoa(highlight.ID))
	}
	params.Set("highlights", strings.Join(highlightList, "-"))

	if p.Folder != "" {
		params.Set("folder_id", p.Folder)
	}

	res, err := svc.Client.Call("/bookmarks/list", params)
	if err != nil {
		return &BookmarkListResponse{}, err
	}
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

// GetText returns the specified bookmark's processed text-view HTML, which is always text/html encoded as UTF-8.
func (svc *BookmarkService) GetText(bookmarkID int) (string, error) {
	params := url.Values{}
	params.Set("bookmark_id", strconv.Itoa(bookmarkID))
	res, err := svc.Client.Call("/bookmarks/get_text", params)
	if err != nil {
		return "", err
	}
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

// Star stars the specified bookmark
func (svc *BookmarkService) Star(bookmarkID int) error {
	params := url.Values{}
	params.Set("bookmark_id", strconv.Itoa(bookmarkID))
	_, err := svc.Client.Call("/bookmarks/star", params)
	return err
}

// UnStar un-stars the specified bookmark
func (svc *BookmarkService) UnStar(bookmarkID int) error {
	params := url.Values{}
	params.Set("bookmark_id", strconv.Itoa(bookmarkID))
	_, err := svc.Client.Call("/bookmarks/unstar", params)
	return err
}

// Archive archives the specified bookmark
func (svc *BookmarkService) Archive(bookmarkID int) error {
	params := url.Values{}
	params.Set("bookmark_id", strconv.Itoa(bookmarkID))
	_, err := svc.Client.Call("/bookmarks/archive", params)
	return err
}

// UnArchive un-archives the specified bookmark
func (svc *BookmarkService) UnArchive(bookmarkID int) error {
	params := url.Values{}
	params.Set("bookmark_id", strconv.Itoa(bookmarkID))
	_, err := svc.Client.Call("/bookmarks/unarchive", params)
	return err
}

// DeletePermanently PERMANENTLY deletes the specified bookmark
func (svc *BookmarkService) DeletePermanently(bookmarkID int) error {
	params := url.Values{}
	params.Set("bookmark_id", strconv.Itoa(bookmarkID))
	_, err := svc.Client.Call("/bookmarks/delete", params)
	return err
}

// Move moves the specified bookmark to the specified folder
func (svc *BookmarkService) Move(bookmarkID int, folderID string) error {
	params := url.Values{}
	params.Set("bookmark_id", strconv.Itoa(bookmarkID))
	params.Set("folder_id", folderID)
	_, err := svc.Client.Call("/bookmarks/move", params)
	return err
}

// UpdateReadProgress updates the read progress on the bookmark
// progress is between 0.0 and 1.0 - a percentage
// when - Unix timestamp - optionally specify when the update happened. If it's set to 0 the current timestamp is used.
func (svc *BookmarkService) UpdateReadProgress(bookmarkID int, progress float32, when int64) error {
	if when == 0 {
		when = time.Now().Unix()
	}
	params := url.Values{}
	params.Set("bookmark_id", strconv.Itoa(bookmarkID))
	params.Set("progress_timestamp", strconv.FormatInt(when, 10))
	params.Set("progress", fmt.Sprintf("%f", progress))
	_, err := svc.Client.Call("/bookmarks/update_read_progress", params)
	return err
}

// Add adds a new bookmark from the specified URL
func (svc *BookmarkService) Add(p BookmarkAddRequestParams) (*Bookmark, error) {
	params := url.Values{}
	params.Set("url", p.URL)
	if p.Description != "" {
		params.Set("description", p.Description)
	}
	if p.Title != "" {
		params.Set("title", p.Title)
	}
	if p.Folder != "" {
		params.Set("folder_id", p.Folder)
	}
	if !p.ResolveFinalURL {
		params.Set("resolve_final_url", "0")
	}
	if p.Content != "" {
		params.Set("content", p.Content)
	}
	if p.PrivateSourceName != "" {
		params.Set("is_private_from_source", p.PrivateSourceName)
	}
	res, err := svc.Client.Call("/bookmarks/add", params)
	if err != nil {
		return nil, err
	}
	var bookmark []Bookmark
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, &APIError{
			StatusCode:   res.StatusCode,
			Message:      err.Error(),
			ErrorCode:    ErrHTTPError,
			WrappedError: err,
		}
	}
	err = json.Unmarshal(bodyBytes, &bookmark)
	if err != nil {
		return nil, &APIError{
			StatusCode:   res.StatusCode,
			Message:      err.Error(),
			ErrorCode:    ErrUnmarshalError,
			WrappedError: err,
		}
	}
	return &bookmark[0], nil
}
