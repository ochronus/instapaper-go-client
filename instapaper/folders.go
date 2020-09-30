package instapaper

import (
	"encoding/json"

	"github.com/gomodule/oauth1/oauth"
)

type Folder struct {
	ID           json.Number `json:"folder_id"`
	Title        string
	Slug         string
	DisplayTitle string `json:"display_title"`
	SyncToMobile int    `json:"sync_to_mobile"`
	Position     int
}

const FOLDER_ID_UNREAD = "unread"
const FOLDER_ID_STARRED = "starred"
const FOLDER_ID_ARCHIVE = "archive"

type FolderService interface {
	List() ([]Folder, error)
}

type FolderServiceOp struct {
	Client      oauth.Client
	Credentials *oauth.Credentials
}

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
