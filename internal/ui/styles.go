package ui

import "github.com/charmbracelet/lipgloss"

var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			PaddingLeft(1)

	windowStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("255"))

	windowDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250"))

	paneStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	paneDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	metaStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	metaBranchStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("72"))

	agentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("141"))

	activeIndicator = lipgloss.NewStyle().
			Foreground(lipgloss.Color("75")).
			Render("●")

	cursorAccent = lipgloss.NewStyle().
			Foreground(lipgloss.Color("75")).
			Render("▎")

	bandA = lipgloss.Color("234")
	bandB = lipgloss.Color("236")

	cursorBg = lipgloss.Color("238")
	cursorFg = lipgloss.Color("255")
)
