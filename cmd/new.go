package cmd

import (
	"fmt"
	"strings"

	"charm.land/huh/v2"
	"github.com/hettiger/clg/cmd/output"
	"github.com/hettiger/clg/internal/changelog"
	"github.com/hettiger/clg/internal/support"
	"github.com/hettiger/clg/internal/validation"
	"github.com/spf13/cobra"
)

type newCmdState struct {
	changeGroup string
	changeType  string
	message     string
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

	if len(app.config.Groups) > 0 {
		groupKeys := support.SortedMapKeys(app.config.Groups)
		newCmd.Flags().StringVarP(
			&state.changeGroup,
			"group",
			"g",
			"",
			fmt.Sprintf("group of change (%s)", strings.Join(groupKeys, ", ")),
		)
	}

	typeKeys := support.SortedMapKeys(app.config.Types)
	newCmd.Flags().StringVarP(
		&state.changeType,
		"type",
		"t",
		"",
		fmt.Sprintf("type of change (%s)", strings.Join(typeKeys, ", ")),
	)

	newCmd.Flags().StringVarP(&state.message, "message", "m", "", "changelog entry")

	return newCmd
}

func addChangelogEntry(app *App, cmd *cobra.Command, state *newCmdState) error {
	groupKeys := support.SortedMapKeys(app.config.Groups)
	typeKeys := support.SortedMapKeys(app.config.Types)

	validateChangeGroup := func(value string) error {
		if len(groupKeys) == 0 {
			return nil
		}

		return validation.ValidateIn("Group of change", value, groupKeys...)
	}

	validateChangeType := func(value string) error {
		return validation.ValidateIn("Type of change", value, typeKeys...)
	}

	var groups []*huh.Group

	if len(groupKeys) > 0 && state.changeGroup == "" {
		options := make([]huh.Option[string], 0, len(groupKeys))
		for _, groupKey := range groupKeys {
			options = append(options, huh.NewOption(app.config.Groups[groupKey], groupKey))
		}

		groups = append(groups, huh.NewGroup(
			huh.NewSelect[string]().
				Title("Group of Change").
				Options(options...).
				Validate(validateChangeGroup).
				Value(&state.changeGroup),
		))
	} else if err := validateChangeGroup(state.changeGroup); err != nil {
		return err
	}

	if state.changeType == "" {
		options := make([]huh.Option[string], 0, len(typeKeys))
		for _, typeKey := range typeKeys {
			options = append(options, huh.NewOption(app.config.Types[typeKey], typeKey))
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
		Group: state.changeGroup,
	}
	path, err := app.entryStore.Write(changelogEntry)
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

func validateTrimmedMessage(value string) error {
	return validation.ValidateMin("Changelog entry", strings.TrimSpace(value), 1)
}
