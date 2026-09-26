package cmd

import (
	"strconv"

	"github.com/hettiger/clg/cmd/output"
	"github.com/spf13/cobra"
)

func NewShowCmd(app *App) *cobra.Command {
	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show unreleased changelog entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			return showUnreleasedChangelogEntries(app)
		},
	}

	return showCmd
}

func showUnreleasedChangelogEntries(app *App) error {
	unreleasedEntries, err := app.entryStore.UnreleasedEntries()
	if err != nil {
		return err
	}

	if len(unreleasedEntries) < 1 {
		return output.PrintSuccess("No changelog entries. Nothing to show.")
	}

	headers := []string{"No."}
	if len(app.config.Groups) > 0 {
		headers = append(headers, "Group")
	}
	headers = append(headers, "Type", "Title", "Branch")
	rows := make([][]string, len(unreleasedEntries))
	for i, entry := range unreleasedEntries {
		rows[i] = []string{strconv.Itoa(i + 1)}
		if len(app.config.Groups) > 0 {
			rows[i] = append(rows[i], entry.Group)
		}
		rows[i] = append(rows[i], entry.Type, entry.Title, entry.Branch)
	}

	return output.PrintTable(headers, rows)
}
