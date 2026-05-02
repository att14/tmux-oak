package ui

import "github.com/charmbracelet/lipgloss"

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("51"))

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("238"))

	windowStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("255"))

	windowDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250"))

	paneActiveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	paneDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246"))

	connectorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("238"))

	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("237")).
			Foreground(lipgloss.Color("255"))

	metaStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243"))

	metaBranchStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("114"))

	agentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("213"))

	activeColor = lipgloss.Color("51")
)
