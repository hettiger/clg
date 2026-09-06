package cmd

import (
	"os"

	"github.com/hettiger/clg/cmd/output"
	"github.com/hettiger/clg/internal/changelog"
	"github.com/spf13/cobra"
)

var (
	marker string
)

var releaseCmd = &cobra.Command{
	Use:   "release [tag]",
	Short: "Build new changelog from unreleased changelog entries",
	Args:  cobra.ExactArgs(1),
	RunE:  addRelease,
}

func init() {
	rootCmd.AddCommand(releaseCmd)

	releaseCmd.Flags().StringVarP(&marker, "marker", "m", "<!-- CLG -->", "insertion marker for new releases")
}

func addRelease(cmd *cobra.Command, args []string) error {
	entryFiles, err := changelog.UnreleasedEntryFiles()
	if err != nil {
		return err
	}

	if len(entryFiles) < 1 {
		return output.PrintSuccess("No logs. Nothing to release.")
	}

	release, err := changelog.ReleaseFromUnreleasedEntries(args[0])
	if err != nil {
		return err
	}

	changelogFile, err := changelog.FileFromCwd(marker)
	if err != nil {
		return err
	}

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
