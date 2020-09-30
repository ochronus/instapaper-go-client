package instapaper

// Highlight represents a highlight within a bookmark
type Highlight struct {
	ID         int `json:"highlight_id"`
	BookmarkID int `json:"bookmark_id"`
	Text       string
	Note       string
	Time       int
	Position   int
}
