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
	go func(db *entities.DB) {
		for t := range ticker.C {
			fmt.Println("Syncing at", t)
			feeds, err := db.GetFeeds()
			if err != nil {
				fmt.Println(err)
				continue
			}
			pubs, err := db.GetPubs()
			if err != nil {
				fmt.Println(err)
				continue
			}
			self, err := db.GetSelf()
			if err != nil {
				fmt.Println(err)
				continue
			}
			if self == nil {
				fmt.Println("No self URL found; skipping a sync")
				continue
			}

			var other []feed.Pub
			for _, p := range pubs {
				if !bytes.Equal(p.URLHash(), self.URLHash()) {
					other = append(other, p)
				}
			}

			newPubs, newFeeds, err := feed.Sync(other, feeds)
			if err != nil {
				fmt.Println(err)
				continue
			}
			for _, f := range newFeeds {
				err = db.PutFeed(f)
				if err != nil {
					fmt.Println(err)
				}
			}
			// update old pubs too to track failures and backoff
			for _, p := range other {
				err = db.PutPub(&p)
				if err != nil {
					fmt.Println(err)
				}
			}
			for _, p := range newPubs {
				err = db.PutPub(&p)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}(db)
	return nil
}
