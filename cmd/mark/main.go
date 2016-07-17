package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/awans/mark/app"
	"github.com/awans/mark/entities"
	"github.com/awans/mark/feed"
	"github.com/awans/mark/server"
	"github.com/docopt/docopt-go"
)

const usage = `mark

Usage:
  mark init
	mark serve`

func initFeed() error {
	markDir := os.Getenv("MARK_DIR")

	if markDir == "" {
		return errors.New("Set the environment variable MARK_DIR")
	}

	err := os.MkdirAll(markDir, 0777)
	if err != nil {
		return err
	}

	store, err := entities.CreateStore(markDir)
	if err != nil {
		return err
	}
	defer store.Close()

	key, err := feed.CreateKeys(markDir)
	if err != nil {
		return err
	}

	feed, err := feed.New(key)
	if err != nil {
		return err
	}

	fp, err := feed.Fingerprint()
	if err != nil {
		return err
	}
	db := entities.NewDB(store, fp)
	return db.PutFeed(feed)
}

func openDbAndKeys() (*rsa.PrivateKey, *entities.DB, error) {
	markDir := os.Getenv("MARK_DIR")
	if markDir == "" {
		return nil, nil, errors.New("Set the environment variable MARK_DIR")
	}

	key, err := feed.OpenKeys(markDir)
	if err != nil {
		return nil, nil, err
	}
	store, err := entities.OpenStore(markDir)
	if err != nil {
		return nil, nil, err
	}

	fp, err := feed.Fingerprint(&key.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	db := entities.NewDB(store, fp)
	db.RebuildIndexes()

	return key, db, nil
}

func serve(db *entities.DB, key *rsa.PrivateKey) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		db.Close()
		os.Exit(0)
	}()

	appDB := app.NewDB(db, key)

	s := server.New(appDB)
	fmt.Printf("Now serving on :8081\n")
	return http.ListenAndServe(":8081", s)
}

func main() {
	args, _ := docopt.Parse(usage, nil, true, "Mark 0", false)

	if args["init"].(bool) {
		err := initFeed()
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
		if args["serve"].(bool) {
			err = serve(db, key)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
