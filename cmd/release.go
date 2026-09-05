package cmd

import (
	"time"

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

	cmd.Printf("## [%s] - %s", release.Tag, time.Now().UTC().Format("2006-01-02"))
	cmd.Println()

	for _, groupType := range changelog.SupportedTypes() {
		if groupType.Keyword == "ignore" {
			continue
		}

		groupedEntries := release.Groups[groupType]
		if len(groupedEntries) == 0 {
			continue
		}

		groupCount := len(groupedEntries)
		groupCountSuffix := "change"
		if groupCount > 1 {
			groupCountSuffix = "changes"
		}

		cmd.Println()
		cmd.Printf("### %s (%v %s)", groupType.Label, groupCount, groupCountSuffix)
		cmd.Println()
		cmd.Println()

		for _, groupEntry := range groupedEntries {
			cmd.Println("- " + groupEntry.Title)
		}
	}

	return nil
}
