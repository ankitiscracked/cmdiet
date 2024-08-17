package ui

import "github.com/charmbracelet/lipgloss"

func TerminalInfo(text string) string {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#32de84")).
		Padding(2).
		Render(lipgloss.NewStyle().Foreground(lipgloss.Color("#874BFC")).Bold(true).Render(text))
}
