package mark

// Bookmark is a model class representing a bookmark
type Bookmark struct {
	ID  string `json:"id"`
	URL  string `json:"url"`
	Note string `json:"note"`
}
