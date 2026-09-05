package cmd

import (
	"strings"

	"charm.land/huh/v2"
	"github.com/hettiger/clg/internal/changelog"
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
		options := make([]huh.Option[string], len(changelog.SupportedTypes()))
		for i, t := range changelog.SupportedTypes() {
			options[i] = huh.NewOption(t.Label, t.Keyword)
		}

		groups = append(groups, huh.NewGroup(
			huh.NewSelect[string]().
				Title("Type of Change").
				Options(options...).
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
		keymap := huh.NewDefaultKeyMap()
		keymap.Input.Prev.Unbind()
		keymap.Select.Prev.Unbind()

		if err := huh.NewForm(groups...).WithKeyMap(keymap).Run(); err != nil {
			return err
		}
	}

	message = strings.TrimSpace(message)

	entry := changelog.Entry{
		Title: message,
		Type:  changeType,
	}
	err := entry.WriteToFile()
	if err != nil {
		return err
	}

	cmd.Println("Success!")

	return nil
}

func validateChangeType(value string) error {
	return validation.ValidateIn(
		"Type of change",
		value,
		changelog.SupportedTypeKeywords()...,
	)
}

func validateTrimmedMessage(value string) error {
	return validation.ValidateMin("Changelog entry", strings.TrimSpace(value), 1)
}
