package cmd

import (
	"charm.land/huh/v2"
	"github.com/spf13/cobra"
)

var (
	changeType string
	message    string
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Add a new changelog entry",
	Args:  cobra.NoArgs,
	RunE:  addChangelogEntry,
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&changeType, "type", "t", "", "type of change")
	newCmd.Flags().StringVarP(&message, "message", "m", "", "changelog entry")
}

func addChangelogEntry(cmd *cobra.Command, args []string) error {
	var groups []*huh.Group

	if changeType == "" {
		groups = append(groups, huh.NewGroup(
			huh.NewSelect[string]().
				Title("Type of Change").
				Options(
					huh.NewOption("New Feature", "added"),
					huh.NewOption("Bug Fix", "fixed"),
					huh.NewOption("Hotfix", "hotfix"),
					huh.NewOption("Feature Change", "changed"),
					huh.NewOption("New Deprecation", "deprecated"),
					huh.NewOption("Feature Removal", "removed"),
					huh.NewOption("Security Fix", "security"),
					huh.NewOption("Performance Improvement", "performance"),
					huh.NewOption("Other", "other"),
					huh.NewOption("No Changelog", "ignore"),
				).
				Value(&changeType),
		))
	}

	if message == "" {
		groups = append(groups, huh.NewGroup(
			huh.NewText().
				Title("Changelog Entry").
				Value(&message),
		))
	}

	if len(groups) > 0 {
		if err := huh.NewForm(groups...).Run(); err != nil {
			return err
		}
	}

	cmd.Println(changeType)
	cmd.Println(message)

	return nil
}
