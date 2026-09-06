package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCommand(app *App) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "clg",
		Short: "Maintain CHANGELOG.md files with ease.",
	}

	rootCmd.AddCommand(
		NewCleanCmd(app),
		NewNewCmd(app),
		NewReleaseCmd(app),
		NewShowCmd(app),
	)

	return rootCmd
}
