package cmd

import (
	"strings"

	"charm.land/huh/v2"
	"github.com/hettiger/clg/internal/validation"
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
				Validate(validateChangeType).
				Value(&changeType),
		))
	} else if err := validateChangeType(changeType); err != nil {
		return err
	}

	if message == "" {
		groups = append(groups, huh.NewGroup(
			huh.NewInput().
				Title("Changelog Entry").
				Validate(validateTrimmedMessage).
				Value(&message),
		))
	} else if err := validateTrimmedMessage(message); err != nil {
		return err
	}

	if len(groups) > 0 {
		if err := huh.NewForm(groups...).Run(); err != nil {
			return err
		}
	}

	message = strings.TrimSpace(message)

	cmd.Println(changeType)
	cmd.Println(message)

	return nil
}

func validateChangeType(value string) error {
	return validation.ValidateIn(
		"Type of change",
		value,
		"added",
		"fixed",
		"hotfix",
		"changed",
		"deprecated",
		"removed",
		"security",
		"performance",
		"other",
		"ignore",
	)
}

func validateTrimmedMessage(value string) error {
	return validation.ValidateMin("Changelog entry", strings.TrimSpace(value), 1)
}
