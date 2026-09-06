package cmd

import (
	"os"

	"charm.land/huh/v2"
	"github.com/hettiger/clg/cmd/output"
	"github.com/spf13/cobra"
)

type cleanCmdState struct {
	isConfirmed bool
}

func NewCleanCmd(app *App) *cobra.Command {
	state := &cleanCmdState{}

	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "Remove all unreleased changelog entries",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeUnreleasedChangelogEntries(app, state)
		},
	}

	cleanCmd.Flags().BoolVarP(&state.isConfirmed, "force", "f", false, "Force remove files")

	return cleanCmd
}

func removeUnreleasedChangelogEntries(app *App, state *cleanCmdState) error {
	files, err := app.changelogEntryStore.UnreleasedEntryFiles()
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return output.PrintSuccess("No logs. Nothing to delete.")
	}

	if !state.isConfirmed {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Delete all unreleased changelog entry files?").
					Value(&state.isConfirmed),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}
	}

	if !state.isConfirmed {
		return nil
	}

	for _, f := range files {
		if err := os.Remove(f.Path); err != nil {
			return err
		}
	}

	return output.PrintSuccess("Done")
}
