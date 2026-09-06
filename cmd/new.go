package cmd

import (
	"strings"

	"charm.land/huh/v2"
	"github.com/hettiger/clg/cmd/output"
	"github.com/hettiger/clg/internal/changelog"
	"github.com/hettiger/clg/internal/validation"
	"github.com/spf13/cobra"
)

type newCmdState struct {
	changeType string
	message    string
}

func NewNewCmd(app *App) *cobra.Command {
	state := &newCmdState{}

	newCmd := &cobra.Command{
		Use:   "new",
		Short: "Add a new changelog entry",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addChangelogEntry(app, cmd, state)
		},
	}

	newCmd.Flags().StringVarP(&state.changeType, "type", "t", "", "type of change")
	newCmd.Flags().StringVarP(&state.message, "message", "m", "", "changelog entry")

	return newCmd
}

func addChangelogEntry(app *App, cmd *cobra.Command, state *newCmdState) error {
	var groups []*huh.Group

	if state.changeType == "" {
		options := make([]huh.Option[string], len(changelog.SupportedTypes()))
		for i, t := range changelog.SupportedTypes() {
			options[i] = huh.NewOption(t.Label, t.Keyword)
		}

		groups = append(groups, huh.NewGroup(
			huh.NewSelect[string]().
				Title("Type of Change").
				Options(options...).
				Validate(validateChangeType).
				Value(&state.changeType),
		))
	} else if err := validateChangeType(state.changeType); err != nil {
		return err
	}

	if state.message == "" {
		groups = append(groups, huh.NewGroup(
			huh.NewInput().
				Title("Changelog Entry").
				Validate(validateTrimmedMessage).
				Value(&state.message),
		))
	} else if err := validateTrimmedMessage(state.message); err != nil {
		return err
	}

	if len(groups) > 0 {
		keymap := huh.NewDefaultKeyMap()
		keymap.Input.Prev.Unbind()
		keymap.Select.Prev.Unbind()

		form := huh.NewForm(groups...).
			WithKeyMap(keymap)

		if err := form.Run(); err != nil {
			return err
		}
	}

	state.message = strings.TrimSpace(state.message)

	changelogEntry := changelog.ChangelogEntry{
		Title: state.message,
		Type:  state.changeType,
	}
	path, err := app.changelogEntryStore.Write(changelogEntry)
	if err != nil {
		return err
	}

	yaml, err := changelogEntry.YAML()
	if err != nil {
		return err
	}

	output.PrintSuccess("Stored changelog entry at: " + path)

	cmd.Print(yaml)

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
