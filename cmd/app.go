package cmd

import (
	"time"

	"github.com/hettiger/clg/internal/changelog"
)

type App struct {
	now                 func() time.Time
	rootDir             string
	changelogEntryStore changelog.EntryStore
}

func NewApp(now func() time.Time, root string, entryStore changelog.EntryStore) *App {
	return &App{now: now, rootDir: root, changelogEntryStore: entryStore}
}

func (a *App) Execute() error {
	return NewRootCommand(a).Execute()
}
