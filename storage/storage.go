package storage

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger/v4"
)

var Storage *badger.DB

func Initialize() {
	location := os.Getenv("TOMATO_STORAGE")
	if location == "" {
		location = "/tmp/tomato/"
	}

	config := badger.DefaultOptions(location)
	config.Logger = nil

	var err error
	Storage, err = badger.Open(config)
	if err != nil {
		fmt.Printf("Cannot open storage: %s\n", err.Error())
		os.Exit(1)
	}
}
