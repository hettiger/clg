package cmd

import (
	"github.com/hettiger/clg/internal/changelog"
	"github.com/spf13/cobra"
)

var releaseCmd = &cobra.Command{
	Use:   "release [tag]",
	Short: "Build new changelog from unreleased changelog entries",
	Args:  cobra.ExactArgs(1),
	RunE:  buildNewChangelog,
}

func init() {
	rootCmd.AddCommand(releaseCmd)
}

func buildNewChangelog(cmd *cobra.Command, args []string) error {
	release, err := changelog.ReleaseFromUnreleasedEntries(args[0])
	if err != nil {
		return err
	}

	markdown, err := release.Markdown()
	if err != nil {
		return err
	}

	cmd.Print(markdown)

	return nil
}
