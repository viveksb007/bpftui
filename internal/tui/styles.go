package tui

import "github.com/charmbracelet/lipgloss"

// Style definitions for consistent theming across the TUI.
var (
	// titleStyle is used for view titles and headers.
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	// selectedStyle highlights the currently selected item.
	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Bold(true)

	// helpStyle is used for help text and keyboard shortcuts.
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	// errorStyle is used for error messages.
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	// normalStyle is used for regular text.
	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	// dimStyle is used for less important text.
	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)
