package app

// Bookmark is a model class representing a bookmark
type Bookmark struct {
	ID     string `json:"id"`
	Title  string `json:"title"` // set by the client
	URL    string `json:"url"`
	Note   string `json:"note"`
	UserID string `json:"user_id", edb:"ref"`
}

// User represents a user in the marks system
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	FP   string `json:"fp"` // latest pub key fingerprint
}
