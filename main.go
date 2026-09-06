package main

import (
	"log"
	"os"
	"time"

	"github.com/hettiger/clg/cmd"
	"github.com/hettiger/clg/internal/changelog"
)

func main() {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	now := func() time.Time {
		return time.Now().UTC()
	}

	entryStore := changelog.NewEntryStore(rootDir, now)

	app := cmd.NewApp(now, rootDir, entryStore)

	if err := app.Execute(); err != nil {
		log.Fatal(err)
	}
}
