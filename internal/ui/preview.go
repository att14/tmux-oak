package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	previewHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	previewContentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("242"))
)

const previewLines = 8

func renderPreview(content string, target string, width int, maxLines int) string {
	var sb strings.Builder

	header := fmt.Sprintf("─── %s ", target)
	remaining := width - lipgloss.Width(header)
	if remaining > 0 {
		header += strings.Repeat("─", remaining)
	}
	sb.WriteString(previewHeaderStyle.Render(header))
	sb.WriteByte('\n')

	lines := strings.Split(content, "\n")

	// Take the last maxLines non-empty lines from the bottom
	var visible []string
	for i := len(lines) - 1; i >= 0 && len(visible) < maxLines; i-- {
		line := lines[i]
		if strings.TrimSpace(line) == "" && len(visible) == 0 {
			continue
		}
		visible = append([]string{line}, visible...)
	}

	for _, line := range visible {
		if lipgloss.Width(line) > width {
			line = truncate(line, width)
		}
		sb.WriteString(previewContentStyle.Render(line))
		sb.WriteByte('\n')
	}

	return sb.String()
}

func truncate(s string, maxWidth int) string {
	w := 0
	for i, r := range s {
		w++
		if w >= maxWidth-1 {
			return s[:i] + "…"
		}
		_ = r
	}
	return s
}
