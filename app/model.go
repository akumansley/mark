package app

// Bookmark is a model class representing a bookmark
type Bookmark struct {
	ID        string `json:"id"`
	CreatedAt int    `json:"created_at"`
	FeedID    string `json:"feed_id"`
	Title     string `json:"title"` // set by the client
	URL       string `json:"url"`
	Note      string `json:"note"`
}

// Profile represents a user in the marks system
type Profile struct {
	ID     string `json:"id"`
	FeedID string `json:"feed_id"`
	Name   string `json:"name"`
}
