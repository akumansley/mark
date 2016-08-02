package app

import (
	"errors"
	"time"

	"github.com/awans/mark/entities"
	"github.com/awans/mark/feed"
)

// DB is the application-level DB interface
type DB struct {
	e *entities.DB
}

// NewDB makes a new app db from an entity db
func NewDB(e *entities.DB) *DB {
	return &DB{e: e}
}

// Close closes the underlying db
func (db *DB) Close() {
	db.e.Close()
}

// GetStream returns a user's stream
func (db *DB) GetStream(count, offset int) ([]Bookmark, error) {
	var bookmarks []Bookmark
	q := db.e.NewQuery("Bookmark").Order("-CreatedAt").Limit(count).Offset(offset)
	q.GetAll(&bookmarks)
	return bookmarks, nil
}

// AddBookmark inserts a bookmark into the db
func (db *DB) AddBookmark(b *Bookmark) error {
	b.CreatedAt = int(time.Now().Unix())
	id, err := db.e.Add(b)
	b.ID = id
	return err
}

// GetUserProfile returns the current user's profile
func (db *DB) GetUserProfile() (*Profile, error) {
	feed, err := db.e.UserFeed()
	if err != nil {
		return nil, err
	}

	fp, err := feed.Fingerprint()
	if err != nil {
		return nil, err
	}

	var ps []Profile
	db.e.NewQuery("Profile").Filter("FeedID =", string(fp)).GetAll(&ps)
	var p *Profile
	if len(ps) == 0 {
		p = &Profile{}
		db.e.Add(p)
	} else {
		p = &ps[0]
	}
	return p, nil
}

// GetProfile gets a profile by feedID
func (db *DB) GetProfile(feedID string) (*Profile, error) {
	var ps []Profile
	db.e.NewQuery("Profile").Filter("FeedID =", feedID).GetAll(&ps)

	if len(ps) == 0 {
		return nil, errors.New("Profile not found")
	}
	return &ps[0], nil

}

// SetProfile sets the current user's Profile
func (db *DB) SetProfile(p *Profile) error {
	old, err := db.GetUserProfile()
	if err != nil {
		return err
	}
	return db.e.Put(old.ID, p)
}

// GetPubs returns all pubs
func (db *DB) GetPubs() ([]feed.Pub, error) {
	return db.e.GetPubs()
}

// PutPub adds a pub
func (db *DB) PutPub(p *feed.Pub) error {
	return db.e.PutPub(p)
}

// GetFeeds returns all feeds
func (db *DB) GetFeeds() ([]feed.SignedFeed, error) {
	return db.e.GetFeeds()
}

// GetFeed returns a specific feed
func (db *DB) GetFeed(id string) (feed.SignedFeed, error) {
	return db.e.GetFeed(id)
}

// PutFeed sets a feed in the db
func (db *DB) PutFeed(f feed.SignedFeed) error {
	return db.e.PutFeed(f)
}
