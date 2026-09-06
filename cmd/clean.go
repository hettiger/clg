package cmd

import (
	"fmt"
	"os"

	"charm.land/huh/v2"
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

	cleanCmd.Flags().BoolVarP(&isConfirmed, "force", "f", false, "Force removing files")
}

func removeUnreleasedChangelogEntries(cmd *cobra.Command, args []string) error {
	files, err := changelog.UnreleasedEntryFiles()
	if err != nil {
		return err
	}

	if len(files) == 0 {
		cmd.Println("No logs. Nothing to delete.")

		return nil
	}

	if !isConfirmed {
		for _, f := range files {
			cmd.Println(f.Path)
		}

		cmd.Println()

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title(fmt.Sprintf("Do you want to delete these %d files?", len(files))).
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

	cmd.Println("Done")

	return nil
}
