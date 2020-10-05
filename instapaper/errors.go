package instapaper

import "fmt"

const (
	// General errors:

	ErrRateLimitExceeded    = 1040 // Rate-limit exceeded
	ErrNotPremiumAccount    = 1041 // Premium account required
	ErrApplicationSuspended = 1042 // Application is suspended

	// Bookmark errors:

	ErrFullContentRequired      = 1220 // Domain requires full content to be supplied
	ErrDomainNotSupported       = 1221 // Domain has opted out of Instapaper compatibility
	ErrInvalidURL               = 1240 // Invalid URL specified
	ErrInvalidBookmarkID        = 1241 // Invalid or missing bookmark_id
	ErrInvalidFolderID          = 1242 // Invalid or missing folder_id
	ErrInvalidProgress          = 1243 // Invalid or missing progress
	ErrInvalidProgressTimestamp = 1244 // Invalid or missing progress_timestamp
	ErrSuppliedContentRequired  = 1245 // Private bookmarks require supplied content
	ErrUnexpected               = 1250 // Unexpected error when saving bookmark

	// Folder errors:

	ErrInvalidTitle              = 1250 // Invalid or missing title
	ErrDuplicateFolder           = 1251 // User already has a folder with this title
	ErrCannotAddBookmarkToFolder = 1252 // Cannot add bookmarks to this folder

	// Operational errors:

	ErrGeneric = 1500 // Unexpected service error
	ErrTextGen = 1550 // Error generating text version of this URL

	// Highlight Errors:

	ErrEmptyText          = 1600 // Cannot create highlight with empty text
	ErrDuplicateHighlight = 1601 // Duplicate highlight

	// Not coming from Instapaper

	ErrNotAuthenticated = 666 // The client did not authenticate
	ErrUnmarshalError   = 667 // Cannot unmarshal the response from Instapaper's API
	ErrHTTPError        = 668 // A generic HTTP error
)

// APIError represents an error returned by the Instapaper API - a numeric code and a message
type APIError struct {
	ErrorCode    int `json:"error_code"`
	StatusCode   int
	Message      string
	WrappedError error
}

func (r *APIError) Error() string {
	return fmt.Sprintf("status %d: err #%d - %v", r.StatusCode, r.ErrorCode, r.Message)
}
