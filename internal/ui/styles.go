package ui

import "github.com/charmbracelet/lipgloss"

var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			PaddingLeft(1)

	windowStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("255")).
			PaddingLeft(1)

	windowDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")).
			PaddingLeft(1)

	paneStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			PaddingLeft(3)

	paneDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			PaddingLeft(3)

	metaStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			PaddingLeft(3)

	metaBranchStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("72"))

	agentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("141"))

	cursorBar = lipgloss.NewStyle().
			Foreground(lipgloss.Color("75")).
			Render("▎")

	activeIndicator = lipgloss.NewStyle().
			Foreground(lipgloss.Color("75")).
			Render("●")

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("236"))
)
