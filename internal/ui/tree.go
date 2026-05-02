package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/att14/tmux-oak/internal/config"
	"github.com/att14/tmux-oak/internal/detect"
	"github.com/att14/tmux-oak/internal/git"
	"github.com/att14/tmux-oak/internal/tmux"
	"github.com/charmbracelet/lipgloss"
)

type NodeKind int

const (
	WindowNode NodeKind = iota
	PaneNode
)

type TreeNode struct {
	Kind        NodeKind
	WindowIndex int
	PaneIndex   int
	Window      *tmux.Window
	Pane        *tmux.Pane
	IsLastChild bool
}

func buildNodes(state *tmux.State, expanded map[int]bool) []TreeNode {
	if state == nil {
		return nil
	}
	var nodes []TreeNode
	for i := range state.Windows {
		w := &state.Windows[i]
		nodes = append(nodes, TreeNode{
			Kind:        WindowNode,
			WindowIndex: w.Index,
			Window:      w,
		})
		if expanded[w.Index] {
			for j := range w.Panes {
				p := &w.Panes[j]
				nodes = append(nodes, TreeNode{
					Kind:        PaneNode,
					WindowIndex: w.Index,
					PaneIndex:   p.Index,
					Pane:        p,
					IsLastChild: j == len(w.Panes)-1,
				})
			}
		}
	}
	return nodes
}

func renderTree(nodes []TreeNode, cursor int, expanded map[int]bool, width int, cfg config.Config, agents map[int]detect.Agent) string {
	var sb strings.Builder

	// Minimal header
	sb.WriteString(headerStyle.Render("oak"))
	sb.WriteByte('\n')
	sb.WriteByte('\n')

	prevWasPane := false
	for i, node := range nodes {
		// Blank line between window groups
		if node.Kind == WindowNode && i > 0 {
			if prevWasPane {
				sb.WriteByte('\n')
			}
		}

		lines := renderNode(node, expanded, width, cfg, agents)
		selected := i == cursor

		for j, line := range lines {
			if j > 0 {
				sb.WriteByte('\n')
			}
			if selected {
				sb.WriteString(cursorBar)
				sb.WriteString(highlightLine(line, width-1))
			} else {
				sb.WriteString(" ")
				sb.WriteString(line)
			}
		}
		sb.WriteByte('\n')

		prevWasPane = node.Kind == PaneNode
	}

	return sb.String()
}

func renderNode(node TreeNode, expanded map[int]bool, width int, cfg config.Config, agents map[int]detect.Agent) []string {
	switch node.Kind {
	case WindowNode:
		return renderWindowNode(node, expanded, width)
	case PaneNode:
		return renderPaneNode(node, width, cfg, agents)
	}
	return nil
}

func renderWindowNode(node TreeNode, expanded map[int]bool, width int) []string {
	w := node.Window

	label := fmt.Sprintf("%d  %s", w.Index, w.Name)

	if w.Active {
		// Right-align the active indicator
		styledLabel := windowStyle.Render(label)
		labelWidth := lipgloss.Width(styledLabel) + 2 // +2 for outer padding
		indicatorWidth := lipgloss.Width(activeIndicator)
		pad := width - labelWidth - indicatorWidth
		if pad < 1 {
			pad = 1
		}
		return []string{styledLabel + strings.Repeat(" ", pad) + activeIndicator}
	}

	return []string{windowDimStyle.Render(label)}
}

func renderPaneNode(node TreeNode, width int, cfg config.Config, agents map[int]detect.Agent) []string {
	p := node.Pane
	var lines []string

	label := p.Command
	if agent, ok := agents[p.PID]; ok {
		label += "  " + agentStyle.Render(agent.Icon)
	}

	if p.Active {
		lines = append(lines, paneStyle.Render(label))
	} else {
		lines = append(lines, paneDimStyle.Render(label))
	}

	if cfg.ShowCwd || cfg.ShowGit {
		var parts []string
		if cfg.ShowCwd {
			parts = append(parts, filepath.Base(p.CurrentPath))
		}
		if cfg.ShowGit {
			if branch := git.Branch(p.CurrentPath); branch != "" {
				parts = append(parts, metaBranchStyle.Render(branch))
			}
		}
		if len(parts) > 0 {
			lines = append(lines, metaStyle.Render(strings.Join(parts, "  ")))
		}
	}

	return lines
}

func highlightLine(s string, width int) string {
	plain := stripAnsi(s)
	padded := plain + strings.Repeat(" ", max(0, width-lipgloss.Width(plain)))
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("236")).
		Render(padded)
}

func stripAnsi(s string) string {
	var result strings.Builder
	inEsc := false
	for _, r := range s {
		if r == '\033' {
			inEsc = true
			continue
		}
		if inEsc {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEsc = false
			}
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}
