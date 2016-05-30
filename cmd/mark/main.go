package main

import (
  "fmt"
  "os"
  "log"
  "github.com/docopt/docopt-go"
  "github.com/awans/mark"
  "errors"
  "crypto/rsa"
)
const usage = `mark

Usage:
  mark init
  mark `

func initDbAndKeys() error {
  markDir := os.Getenv("MARK_DIR")

  if  markDir == "" {
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

func openDbAndKeys() (*rsa.PrivateKey, mark.Store, error) {
  markDir := os.Getenv("MARK_DIR")
  if  markDir == "" {
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

  return key, store, nil
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
    key, store, err := openDbAndKeys()
    fmt.Println(key, store, err)
  }
}
