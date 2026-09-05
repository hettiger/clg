package cmd

import (
	"github.com/spf13/cobra"
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Build new changelog from unreleased changelog entries",
	RunE:  buildNewChangelog,
}

func init() {
	rootCmd.AddCommand(releaseCmd)
}

func buildNewChangelog(cmd *cobra.Command, args []string) error {
	return nil
}
