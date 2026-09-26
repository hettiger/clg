package cmd

import (
	"time"

	"github.com/hettiger/clg/internal/changelog"
	"github.com/hettiger/clg/internal/config"
)

type App struct {
	config     config.Config
	now        func() time.Time
	rootDir    string
	entryStore changelog.EntryStore
	gitService changelog.GitServcie
}

func NewApp(
	config config.Config,
	now func() time.Time,
	root string,
	entryStore changelog.EntryStore,
	gitService changelog.GitServcie,
) *App {
	return &App{
		config:     config,
		now:        now,
		rootDir:    root,
		entryStore: entryStore,
		gitService: gitService,
	}
}

func (a *App) Execute() error {
	return NewRootCommand(a).Execute()
}
