package cmd

import (
	"strconv"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
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

	t := table.New().
		Headers("No.", "Type", "Log", "Author").
		Border(lipgloss.ASCIIBorder()).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return lipgloss.NewStyle().
					Foreground(lipgloss.Green).
					Padding(0, 1)
			default:
				return lipgloss.NewStyle().
					Padding(0, 1)
			}
		})

	for i, entry := range unreleasedEntries {
		t.Row(strconv.Itoa(i+1), entry.Type, entry.Title, entry.Author)
	}

	lipgloss.Println(t)

	return nil
}
