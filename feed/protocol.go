package feed

import (
	"fmt"
	"log"
	"time"
)

// Protocol related paths
const (
	ProtocolRoot = "sync"
	PubsPath     = "pubs"
	HeadsPath    = "heads"
	AnnouncePath = "announce"
	FeedPath     = "feed"
)

type pubLen struct {
	Pub *Pub
	Len int
}

// Sync gets any new updates from the list of pubs.
// It works incrementally on top of the feeds passed in
// so pass in all known feeds and pubs
// TODO: only load the feed delta
func Sync(pubs []Pub, feeds []SignedFeed) ([]Pub, []SignedFeed, error) {
	feedsByID := make(map[string]SignedFeed)

	for _, feed := range feeds {
		fp, err := feed.Fingerprint()
		if err != nil {
			return nil, nil, err
		}
		feedsByID[string(fp)] = feed
	}

	feedPubs := make(map[string]pubLen)
	pubsByURL := make(map[string]Pub)

	for i := range pubs {
		pub := &pubs[i]
		if pub.ShouldUpdate() {
			pub.LastChecked = time.Now().Unix()
			pubsToAdd, err := pub.GetPubs()
			if err != nil {
				log.Println(err)
				pub.Failures++
				continue
			}
			pub.Failures = 0 // reset -- we had a sucessful response, so it's still alive

			for _, pubToAdd := range pubsToAdd {
				key := string(pubToAdd.URLHash())
				if _, ok := pubsByURL[key]; !ok {
					pubsByURL[key] = pubToAdd
				}
			}

			heads, err := pub.GetHeads()
			if err != nil {
				log.Println(err)
				pub.Failures++
				continue
			}

			for _, head := range heads {
				// do we have this feed at all
				if f, ok := feedsByID[head.ID]; ok {
					fp, err := f.Fingerprint()
					if err != nil {
						return nil, nil, err
					}
					// best so far
					best := len(f)
					if pl, ok := feedPubs[string(fp)]; ok {
						best = pl.Len
					}
					// is theirs better
					if head.Len > best {
						feedPubs[string(fp)] = pubLen{Pub: pub, Len: head.Len}
					}
				} else {
					// we didn't have this feed, so add it
					feedPubs[head.ID] = pubLen{Pub: pub, Len: head.Len}
				}
			}
		}
	}

	// now we know where the latest feeds are, so let's get 'em
	var outFeeds []SignedFeed
	for fp, pl := range feedPubs {
		pub := pl.Pub
		pub.LastUpdated = time.Now().Unix()
		feed, err := pub.GetFeed(fp)
		if err != nil {
			log.Println(err)
			continue
		}
		outFeeds = append(outFeeds, *feed)
	}

	// save off any new friends
	var outPubs []Pub
	for _, p := range pubsByURL {
		outPubs = append(outPubs, p)
	}

	return outPubs, outFeeds, nil
}

// Announce tells your known pubs about some update to a feed
func Announce(self *Pub, pubs []Pub, f SignedFeed) error {
	fmt.Printf("Gonna announce myself: %s\n", self.URL)
	fp, err := f.Fingerprint()
	if err != nil {
		panic(err)
	}
	head := Head{ID: string(fp), Len: len(f)}

	a := Announcement{Pub: *self, Heads: []Head{head}}
	for _, p := range pubs {
		p.Announce(&a)
	}
	return nil
}
