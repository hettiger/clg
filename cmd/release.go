package cmd

import (
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
	release, err := changelog.ReleaseFromUnreleasedEntries(args[0])
	if err != nil {
		return err
	}

	file, err := changelog.FileFromCwd(marker)
	if err != nil {
		return err
	}

	return file.AddRelease(release)
}
