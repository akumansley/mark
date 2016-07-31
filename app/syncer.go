package app

import (
	"bytes"
	"fmt"
	"time"

	"github.com/awans/mark/entities"
	"github.com/awans/mark/feed"
)

// Sync starts a goroutine that runs sync every durationSpec
func Sync(durationSpec string, db *entities.DB) error {
	duration, err := time.ParseDuration(durationSpec)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(duration)
	go func() {
		for t := range ticker.C {
			fmt.Println("Syncing at", t)
			feeds, err := db.GetFeeds()
			pubs, err := db.GetPubs()
			self, err := db.GetSelf()

			var other []feed.Pub
			for _, p := range pubs {
				if !bytes.Equal(p.URLHash(), self.URLHash()) {
					other = append(other, p)
				}
			}

			newPubs, newFeeds, err := feed.Sync(other, feeds)
			if err != nil {
				continue
			}
			for _, f := range newFeeds {
				db.PutFeed(f)
			}
			for _, p := range newPubs {
				db.PutPub(&p)
			}
		}
	}()
	return nil
}
