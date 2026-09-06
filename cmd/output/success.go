package output

import "charm.land/lipgloss/v2"

var successStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Green)

func PrintSuccess(message string) error {
	_, err := lipgloss.Println(successStyle.Render(message))
	return err
}
