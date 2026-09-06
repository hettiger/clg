package cmd

import (
	"os"

	"github.com/hettiger/clg/cmd/output"
	"github.com/hettiger/clg/internal/changelog"
	"github.com/spf13/cobra"
)

type releaseCmdState struct {
	marker string
}

func NewReleaseCmd(app *App) *cobra.Command {
	state := &releaseCmdState{}

	releaseCmd := &cobra.Command{
		Use:   "release [tag]",
		Short: "Add new release to the changelog",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addRelease(app, args, state)
		},
	}

	releaseCmd.Flags().StringVarP(
		&state.marker,
		"marker",
		"m",
		"<!-- CLG -->",
		"insertion marker for new releases",
	)

	return releaseCmd
}

func addRelease(app *App, args []string, state *releaseCmdState) error {
	entryFiles, err := app.changelogEntryStore.UnreleasedEntryFiles()
	if err != nil {
		return err
	}

	if len(entryFiles) < 1 {
		return output.PrintSuccess("No changelog entries. Nothing to release.")
	}

	unreleasedEntries, err := app.changelogEntryStore.UnreleasedEntries()
	if err != nil {
		return err
	}

	release, err := changelog.NewRelease(args[0], unreleasedEntries, app.now())
	if err != nil {
		return err
	}

	changelogFile := changelog.NewChangelogFile(app.rootDir, state.marker)

	markdown, err := changelogFile.AddRelease(release)
	if err != nil {
		return err
	}

	for _, f := range entryFiles {
		if err := os.Remove(f.Path); err != nil {
			return err
		}
	}

	return output.PrintSuccess(markdown)
}
