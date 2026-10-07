package main

import (
	"log"
	"os"
	"time"
	_ "time/tzdata"
	"uuid"

	"github.com/hettiger/clg/cmd"
	"github.com/hettiger/clg/internal/changelog"
	"github.com/hettiger/clg/internal/config"
)

func main() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.Load(userHomeDir, rootDir)
	if err != nil {
		log.Fatal(err)
	}

	now := func() time.Time {
		return time.Now().UTC()
	}

	uuidV7 := func() string {
		return uuid.NewV7().String()
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
