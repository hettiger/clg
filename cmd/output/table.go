package output

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

var headerRowStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Green).
	Padding(0, 1)

var defaultStyle = lipgloss.NewStyle().
	Padding(0, 1)

func PrintTable(headers []string, rows [][]string) error {
	t := table.New().
		Headers(headers...).
		Border(lipgloss.ASCIIBorder()).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerRowStyle
			}

			return defaultStyle
		})

	for _, row := range rows {
		t.Row(row...)
	}

	_, err := lipgloss.Println(t)

	return err
}
