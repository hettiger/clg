package cmd

import (
	"strconv"

	"github.com/hettiger/clg/cmd/output"
	"github.com/hettiger/clg/internal/changelog"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show unreleased changelog entries",
	RunE:  showUnreleasedChangelogEntries,
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func showUnreleasedChangelogEntries(cmd *cobra.Command, args []string) error {
	unreleasedEntries, err := changelog.UnreleasedEntries()
	if err != nil {
		return err
	}

	if len(unreleasedEntries) < 1 {
		return output.PrintSuccess("No logs. Nothing to show.")
	}

	headers := []string{"No.", "Type", "Log", "Author"}
	rows := make([][]string, len(unreleasedEntries))
	for i, entry := range unreleasedEntries {
		rows[i] = []string{strconv.Itoa(i + 1), entry.Type, entry.Title, entry.Author}
	}

	return output.PrintTable(headers, rows)
}
