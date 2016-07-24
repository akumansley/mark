package main

import (
	"crypto/rsa"
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
  mark init [-d <dir>]
	mark serve [-d <dir>] [-p <port>]

Options:
	-d <dir>, --data-dir <dir>  Specify data directory [default: /var/opt/mark]
	-p <port>, --port <port>		Specify port [default: 8080]

`

func initFeed(markDir string) error {
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

func openDbAndKeys(markDir string) (*rsa.PrivateKey, *entities.DB, error) {
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

func serve(db *entities.DB, key *rsa.PrivateKey, dir, port string) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		db.Close()
		os.Exit(0)
	}()

	appDB := app.NewDB(db, key)

	s := server.New(appDB)
	fmt.Printf("Now serving on :%s\n", port)
	return http.ListenAndServe(":"+port, s)
}

func main() {
	args, _ := docopt.Parse(usage, nil, true, "Mark 0", false)
	dir := args["--data-dir"].(string)

	if args["init"].(bool) {
		err := initFeed(dir)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	} else {
		key, db, err := openDbAndKeys(dir) // maybe wrap this in a Session
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		if args["serve"].(bool) {
			port := args["--port"].(string)
			err = serve(db, key, dir, port)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
