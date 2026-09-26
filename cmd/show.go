package cmd

import (
	"strconv"

	"github.com/hettiger/clg/cmd/output"
	"github.com/hettiger/clg/internal/changelog"
	"github.com/spf13/cobra"
)

type showCmdState struct {
	branch string
}

func NewShowCmd(app *App) *cobra.Command {
	state := &showCmdState{}

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show unreleased changelog entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			return showUnreleasedChangelogEntries(app, state)
		},
	}

	showCmd.Flags().StringVarP(
		&state.branch,
		"branch",
		"b",
		"",
		"filter entries by branch",
	)

	return showCmd
}

func showUnreleasedChangelogEntries(app *App, state *showCmdState) error {
	unreleasedEntries, err := app.entryStore.UnreleasedEntries()
	if err != nil {
		return err
	}

	var filteredEntries []changelog.ChangelogEntry
	if state.branch == "" {
		filteredEntries = unreleasedEntries
	} else {
		filteredEntries = make([]changelog.ChangelogEntry, 0, len(unreleasedEntries))
		for _, entry := range unreleasedEntries {
			if entry.Branch == state.branch {
				filteredEntries = append(filteredEntries, entry)
			}
		}
	}

	if len(filteredEntries) < 1 {
		return output.PrintSuccess("No changelog entries. Nothing to show.")
	}

	headers := []string{"No."}
	if len(app.config.Groups) > 0 {
		headers = append(headers, "Group")
	}
	headers = append(headers, "Type", "Title", "Branch")
	rows := make([][]string, len(filteredEntries))
	for i, entry := range filteredEntries {
		rows[i] = []string{strconv.Itoa(i + 1)}
		if len(app.config.Groups) > 0 {
			rows[i] = append(rows[i], entry.Group)
		}
		rows[i] = append(rows[i], entry.Type, entry.Title, entry.Branch)
	}

	return output.PrintTable(headers, rows)
}
