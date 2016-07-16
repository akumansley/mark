package main

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/awans/mark"
	"github.com/awans/mark/server"
	"github.com/davecgh/go-spew/spew"
	"github.com/docopt/docopt-go"
)

const usage = `mark

Usage:
  mark init
  mark list
  mark add <url>
	mark feed
	mark serve`

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
	db.Register(&mark.Bookmark{})
	var bookmarks []mark.Bookmark
	db.GetAll(&key.PublicKey, &bookmarks)

	bookmark := mark.Bookmark{URL: url, Note: ""}
	entity := mark.Entity{ID: strconv.Itoa(len(bookmarks)), Body: &bookmark}
	err := db.Add(key, entity)
	if err != nil {
		return err
	}
	return nil
}

func list(db *mark.DB, key *rsa.PrivateKey) error {
	db.Register(&mark.Bookmark{})
	var bookmarks []mark.Bookmark
	db.GetAll(&key.PublicKey, &bookmarks)

	for i, v := range bookmarks {
		fmt.Printf("%d. %s\n", i+1, v.URL)
	}

	return nil
}

func feed(db *mark.DB, key *rsa.PrivateKey) error {
	feed, err := db.FeedForPubKey(&key.PublicKey)
	if err != nil {
		return err
	}
	fmt.Println("Feed:")
	spew.Dump(feed)

	fmt.Println("\n\nSerialized:")
	bytes, err := json.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", bytes)
	fmt.Println("Current Key:")
	jwk, err := feed.CurrentKey()
	if err != nil {
		return err
	}
	bytes, err = json.MarshalIndent(jwk, "", "  ")
	fmt.Printf("%s", bytes)
	return nil
}

func serve(db *mark.DB, key *rsa.PrivateKey) error {
	db.Register(&mark.Bookmark{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		db.Close()
		os.Exit(0)
	}()
	s := server.New(db, key)
	fmt.Printf("Now serving on :8081\n")
	return http.ListenAndServe(":8081", s)
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
		key, db, err := openDbAndKeys() // maybe wrap this in a Session
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
		if args["feed"].(bool) {
			err = feed(db, key)
			if err != nil {
				log.Fatal(err)
			}
		}
		if args["serve"].(bool) {
			err = serve(db, key)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
