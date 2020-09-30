package instapaper

import (
	"encoding/json"

	"github.com/gomodule/oauth1/oauth"
)

// Folder represents a folder on Instapaper - there are 3 default ones, see FolderIDUnread, FolderIDStarred and FolderIDArchive
type Folder struct {
	ID           json.Number `json:"folder_id"`
	Title        string
	Slug         string
	DisplayTitle string `json:"display_title"`
	SyncToMobile int    `json:"sync_to_mobile"`
	Position     int
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
	Client      oauth.Client
	Credentials *oauth.Credentials
}

// List returns the list of *custom created* folders. It does not return any of the built in ones!
func (svc *FolderServiceOp) List() ([]Folder, error) {
	res, err := svc.Client.Post(nil, svc.Credentials, "https://www.instapaper.com/api/1.1/folders/list", nil)
	if err == nil {
		var folderList []Folder
		err := json.NewDecoder(res.Body).Decode(&folderList)
		if err != nil {
			return nil, err
		} else {
			return folderList, nil
		}
	} else {
		return nil, err
	}
}
