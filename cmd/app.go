package cmd

import (
	"time"

	"github.com/hettiger/clg/internal/changelog"
	"github.com/hettiger/clg/internal/config"
)

type App struct {
	config              config.Config
	now                 func() time.Time
	rootDir             string
	changelogEntryStore changelog.EntryStore
}

func NewApp(config config.Config, now func() time.Time, root string, entryStore changelog.EntryStore) *App {
	return &App{config: config, now: now, rootDir: root, changelogEntryStore: entryStore}
}

func (a *App) Execute() error {
	return NewRootCommand(a).Execute()
}
