package app

import (
	"fmt"
	"time"

	"github.com/awans/mark/entities"
	"github.com/awans/mark/feed"
)

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
			newPubs, newFeeds, err := feed.Sync(pubs, feeds)
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
