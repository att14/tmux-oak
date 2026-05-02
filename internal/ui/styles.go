package ui

import "github.com/charmbracelet/lipgloss"

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("51")).
			PaddingLeft(1)

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("238"))

	windowStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("255"))

	windowDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	paneActiveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	paneDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236"))

	metaDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("239"))

	metaBranchStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("114"))

	activeColor = lipgloss.Color("51")
)
