package instapaper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
)

// Highlight represents a highlight within a bookmark
type Highlight struct {
	ID         int `json:"highlight_id"`
	BookmarkID int `json:"bookmark_id"`
	Text       string
	Note       string
	Time       json.Number
	Position   int
}

type HighlightService struct {
	Client Client
}

// List fetches all highlights for the specified bookmark
func (svc *HighlightService) List(bookmarkID int) ([]Highlight, error) {
	path := fmt.Sprintf("/bookmarks/%d/highlights", bookmarkID)
	res, err := svc.Client.Call(path, nil)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, &APIError{
			StatusCode:   res.StatusCode,
			Message:      err.Error(),
			ErrorCode:    ErrHTTPError,
			WrappedError: err,
		}
	}
	var highlightList []Highlight
	err = json.Unmarshal(bodyBytes, &highlightList)
	if err != nil {
		return nil, &APIError{
			StatusCode:   res.StatusCode,
			Message:      err.Error(),
			ErrorCode:    ErrUnmarshalError,
			WrappedError: err,
		}
	}
	return highlightList, nil
}

// Delete removes the specified highlight
func (svc *HighlightService) Delete(highlightID int) error {
	path := fmt.Sprintf("/highlights/%d/delete", highlightID)
	_, err := svc.Client.Call(path, nil)
	if err != nil {
		return err
	}
	return nil
}

// Add adds a highlight for the specified bookmark
func (svc *HighlightService) Add(bookmarkID int, text string, position int) (*Highlight, error) {
	path := fmt.Sprintf("/bookmarks/%d/highlight", bookmarkID)
	params := url.Values{}
	params.Set("text", text)
	params.Set("position", strconv.Itoa(position))
	res, err := svc.Client.Call(path, params)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, &APIError{
			StatusCode:   res.StatusCode,
			Message:      err.Error(),
			ErrorCode:    ErrHTTPError,
			WrappedError: err,
		}
	}
	var highlightList []Highlight
	err = json.Unmarshal(bodyBytes, &highlightList)
	if err != nil {
		return nil, &APIError{
			StatusCode:   res.StatusCode,
			Message:      err.Error(),
			ErrorCode:    ErrUnmarshalError,
			WrappedError: err,
		}
	}
	return &highlightList[0], nil
}
