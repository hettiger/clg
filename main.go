package main

import (
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hettiger/clg/cmd"
	"github.com/hettiger/clg/internal/changelog"
	"github.com/hettiger/clg/internal/config"
)

func main() {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.Load(rootDir)
	if err != nil {
		log.Fatal(err)
	}

	now := func() time.Time {
		return time.Now().UTC()
	}

	uuidV7 := func() (string, error) {
		u, err := uuid.NewV7()
		if err != nil {
			return "", err
		}

		return u.String(), nil
	}

	entryStore := changelog.NewEntryStore(
		rootDir,
		uuidV7,
		cfg.Groups,
		cfg.Types,
	)

	gitService := changelog.NewGitService(rootDir)

	app := cmd.NewApp(cfg, now, rootDir, entryStore, gitService)

	if err := app.Execute(); err != nil {
		log.Fatal(err)
	}
}
