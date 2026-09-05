package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/hettiger/clg/internal/changelog"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show unreleased changes",
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

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
	fmt.Fprintln(w, "NO\tTYPE\tLOG\tAUTHOR")

	for i, entry := range unreleasedEntries {
		fmt.Fprintf(w, "%v\t%s\t%s\t%s\n", i+1, entry.Type, entry.Title, entry.Author)
	}

	w.Flush()

	return nil
}
