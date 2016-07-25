package feed

import "log"

const (
	ProtocolRoot = "sync"
	PubsPath     = "pubs"
	HeadsPath    = "heads"
	AnnouncePath = "announce"
	FeedPath     = "feed"
)

type pubLen struct {
	Pub Pub
	Len int
}

// Sync gets any new updates from the list of pubs.
// It works incrementally on top of the feeds passed in
// so pass in all known feeds and pubs
func Sync(pubs []Pub, feeds []SignedFeed) ([]Pub, []SignedFeed, error) {
	feedsById := make(map[string]SignedFeed)

	for _, feed := range feeds {
		fp, err := feed.Fingerprint()
		if err != nil {
			return nil, nil, err
		}
		feedsById[string(fp)] = feed
	}

	feedPubs := make(map[string]pubLen)
	pubsByUrl := make(map[string]Pub)

	for _, pub := range pubs {
		if pub.ShouldUpdate() {
			pubsToAdd, err := pub.GetPubs()

			if err != nil {
				log.Println(err)
			}

			for _, pubToAdd := range pubsToAdd {
				pubsByUrl[string(pubToAdd.URLHash())] = pubToAdd
			}

			// TODO set last_checked
			heads, err := pub.GetHeads()
			if err != nil {
				// we should eventually fail-out the pub
				log.Println(err)
				continue
			}

			for _, head := range heads {
				// do we have this feed at all
				if f, ok := feedsById[head.ID]; ok {
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
		feed, err := pub.GetFeed(fp)
		if err != nil {
			log.Println(err)
			continue
		}
		outFeeds = append(outFeeds, *feed)
	}

	var outPubs []Pub
	for _, p := range pubsByUrl {
		outPubs = append(outPubs, p)
	}

	return outPubs, outFeeds, nil
}

// Announce tells your known pubs about some update to a feed
func Announce(pubs []Pub, feed Feed) error {
	return nil
}
