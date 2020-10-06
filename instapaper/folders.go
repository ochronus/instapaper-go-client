package instapaper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

// Folder represents a folder on Instapaper - there are 3 default ones, see FolderIDUnread, FolderIDStarred and FolderIDArchive
type Folder struct {
	ID           json.Number `json:"folder_id"`
	Title        string
	Slug         string
	DisplayTitle string `json:"display_title"`
	SyncToMobile int    `json:"sync_to_mobile"`
	Position     json.Number
}

// FolderIDUnread is the default folder - unread bookmarks
const FolderIDUnread = "unread"

// FolderIDStarred is a built-in folder for starred bookmarks
const FolderIDStarred = "starred"

// FolderIDArchive is a built-in folder for archived bookmarks
const FolderIDArchive = "archive"

type folderService interface {
	List() ([]Folder, error)
}

// FolderServiceOp encapsulates all folder operations
type FolderServiceOp struct {
	Client Client
}

// List returns the list of *custom created* folders. It does not return any of the built in ones!
func (svc *FolderServiceOp) List() ([]Folder, error) {
	res, err := svc.Client.Call("/folders/list", nil)
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
	var folderList []Folder
	err = json.Unmarshal(bodyBytes, &folderList)
	if err != nil {
		return nil, &APIError{
			StatusCode:   res.StatusCode,
			Message:      err.Error(),
			ErrorCode:    ErrUnmarshalError,
			WrappedError: err,
		}
	}
	fmt.Println(string(bodyBytes))
	return folderList, nil
}

// Add creates a folder and returns with it if there wasn't already one with the same title - in that case it returns an error
func (svc *FolderServiceOp) Add(title string) (*Folder, error) {
	params := url.Values{}
	params.Set("title", title)
	res, err := svc.Client.Call("/folders/add", params)
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
	var folderList []Folder
	err = json.Unmarshal(bodyBytes, &folderList)
	if err != nil {
		return nil, &APIError{
			StatusCode:   res.StatusCode,
			Message:      err.Error(),
			ErrorCode:    ErrUnmarshalError,
			WrappedError: err,
		}
	}
	return &folderList[0], nil
}

// Delete removes a folder and moves all of its bookmark entries to the archive
func (svc *FolderServiceOp) Delete(folderID string) error {
	params := url.Values{}
	params.Set("folder_id", folderID)
	_, err := svc.Client.Call("/folders/delete", params)
	if err != nil {
		return err
	}
	return nil
}

// SetOrder sets the order of the user-created folders.
// Format: folderid1:order1,folderid2:order2,...,folderidN,orderN
// example: 100:1,200:2,300:3
// the order of the pairs in the list does not matter.
// You should include all folders for consistency.
// !!!No errors returned for missing or invalid folders!!!
func (svc *FolderServiceOp) SetOrder(folderOrderlist string) ([]Folder, error) {
	params := url.Values{}
	params.Set("order", folderOrderlist)
	res, err := svc.Client.Call("/folders/set_order", params)
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
	var folderList []Folder
	err = json.Unmarshal(bodyBytes, &folderList)
	if err != nil {
		return nil, &APIError{
			StatusCode:   res.StatusCode,
			Message:      err.Error(),
			ErrorCode:    ErrUnmarshalError,
			WrappedError: err,
		}
	}
	fmt.Println(string(bodyBytes))
	return folderList, nil
}
