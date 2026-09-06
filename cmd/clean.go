package cmd

import (
	"os"

	"charm.land/huh/v2"
	"github.com/hettiger/clg/cmd/output"
	"github.com/hettiger/clg/internal/changelog"
	"github.com/spf13/cobra"
)

var (
	isConfirmed bool
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove all unreleased changelog entries",
	Args:  cobra.NoArgs,
	RunE:  removeUnreleasedChangelogEntries,
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.Flags().BoolVarP(&isConfirmed, "force", "f", false, "Force remove files")
}

func removeUnreleasedChangelogEntries(cmd *cobra.Command, args []string) error {
	files, err := changelog.UnreleasedEntryFiles()
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return output.PrintSuccess("No logs. Nothing to delete.")
	}

	if !isConfirmed {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Delete all unreleased changelog entry files?").
					Value(&isConfirmed),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}
	}

	if !isConfirmed {
		return nil
	}

	for _, f := range files {
		if err := os.Remove(f.Path); err != nil {
			return err
		}
	}

	return output.PrintSuccess("Done")
}
