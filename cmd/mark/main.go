package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/awans/mark/app"
	"github.com/awans/mark/entities"
	"github.com/awans/mark/feed"
	"github.com/awans/mark/sandstorm"
	"github.com/awans/mark/server"
	"github.com/docopt/docopt-go"
)

const bootstrapURL = "https://marko.awans.org"

const usage = `mark

Usage:
  mark init [-d <dir>]
  mark serve [-d <dir>] [-p <port>]
  mark dump [-d <dir>]
  mark rebuild [-d <dir>]

Options:
	-d <dir>, --data-dir <dir>  Specify data directory [default: /var/opt/mark]
	-p <port>, --port <port>		Specify port [default: 8080]

`

func initFeed(markDir string) error {
	err := os.RemoveAll(markDir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(markDir, 0777)
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
	db := entities.NewDB(store, fp, key)
	_, err = db.PutUserFeed(feed)
	return err
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

	db := entities.NewDB(store, fp, key)
	db.RebuildIndexes()

	return key, db, nil
}

func dump(db *entities.DB) {
	fmt.Printf("%s", db.Dump())
}

func rebuild(db *entities.DB) {
	err := db.RebuildUserFeed()
	if err != nil {
		panic(err)
	}
}

// HTTPGetter implements Getter with the net/http package
type HTTPGetter struct{}

// Get implements Getter
func (g HTTPGetter) Get(url string) (*http.Response, error) {
	res, err := http.Get(url)
	return res, err
}

func serve(db *entities.DB, key *rsa.PrivateKey, port string) error {
	bootstrap := feed.Pub{URL: bootstrapURL, LastUpdated: time.Now().Unix(), LastChecked: time.Now().Unix()}
	db.PutPub(&bootstrap)

	// Catch ctrl-c and gracefully exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		db.Close()
		os.Exit(0)
	}()

	appDB := app.NewDB(db)

	// The server bootstraps a sandstorm getter
	// so it has to be called *before* the app.Sync starts..
	// Note that initial sync will block on the first request to hit the server
	s, sessionBus := server.New(appDB)
	g := sandstorm.NewGetter(sessionBus)
	feed.Initialize(g)

	go sandstorm.GetEndpointURL(db, sessionBus)

	sessionBus.Sub()

	app.Sync("10s", db)

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
			err = serve(db, key, port)
			if err != nil {
				log.Fatal(err)
			}
		} else if args["dump"].(bool) {
			dump(db)
		} else if args["rebuild"].(bool) {
			rebuild(db)
		}
	}
}
