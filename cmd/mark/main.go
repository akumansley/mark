package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/awans/mark"
	"github.com/docopt/docopt-go"
	"log"
	"os"
)

const usage = `mark

Usage:
  mark init
  mark list
  mark add <url>`

func initDbAndKeys() error {
	markDir := os.Getenv("MARK_DIR")

	if markDir == "" {
		return errors.New("Set the environment variable MARK_DIR")
	}

	err := os.MkdirAll(markDir, 0777)
	if err != nil {
		return err
	}

	store, err := mark.CreateStore(markDir)
	if err != nil {
		return err
	}
	defer store.Close()

	_, err = mark.CreateKeys(markDir)
	if err != nil {
		return err
	}
	return nil
}

func openDbAndKeys() (*rsa.PrivateKey, *mark.DB, error) {
	markDir := os.Getenv("MARK_DIR")
	if markDir == "" {
		return nil, nil, errors.New("Set the environment variable MARK_DIR")
	}

	key, err := mark.OpenKeys(markDir)
	if err != nil {
		return nil, nil, err
	}
	store, err := mark.OpenStore(markDir)
	if err != nil {
		return nil, nil, err
	}

	db := mark.DBFromStore(store)

	return key, db, nil
}

func add(db *mark.DB, key *rsa.PrivateKey, url string) error {
	bookmark := mark.Bookmark{URL: url, Note: ""}
	entity := mark.Entity{ID: "1", Body: &bookmark}
	err := db.Add(key, entity)
	if err != nil {
		return err
	}
	return nil
}

func list(db *mark.DB, key *rsa.PrivateKey) error {
	db.Register(&mark.Bookmark{})

	// TODO move this inside of DB setup
	feed, err := db.FeedForPubKey(&key.PublicKey)
	if err != nil {
		return err
	}
	var eavs []mark.EAV
	for _, op := range feed.Ops {
		if op.Op == "eav" {
			lsOfLs := op.Body.([]interface{})

			for _, ls := range lsOfLs {
				eav := mark.EAV{}
				eav.BuildFromList(ls.([]interface{}))
				eavs = append(eavs, eav)
			}

		}
	}

	// TODO get all the EAVS before inflating
	entities := db.Inflate(eavs)
	bookmarks := make([]*mark.Bookmark, len(entities))
	for i, v := range entities {
		bookmarks[i] = v.Body.(*mark.Bookmark)
	}
	for i, v := range bookmarks {
		fmt.Printf("%d. %s\n", i+1, v.URL)
	}

	return nil
}

func main() {
	args, _ := docopt.Parse(usage, nil, true, "Mark 0", false)

	if args["init"].(bool) {
		err := initDbAndKeys()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	} else {
		key, db, err := openDbAndKeys()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		if args["add"].(bool) {
			err = add(db, key, args["<url>"].(string))
			if err != nil {
				log.Fatal(err)
			}
		}

		if args["list"].(bool) {
			err = list(db, key)
			if err != nil {
				log.Fatal(err)
			}

		}
	}
}
