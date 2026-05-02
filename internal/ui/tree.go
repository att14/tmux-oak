package ui

import (
	"fmt"
	"os"
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

	sb.WriteString(headerStyle.Render("oak"))
	sb.WriteByte('\n')

	// Assign a group index to each node (increments on each window)
	groupOf := make([]int, len(nodes))
	group := 0
	for i, node := range nodes {
		if node.Kind == WindowNode && i > 0 {
			group++
		}
		groupOf[i] = group
	}

	for i, node := range nodes {
		bg := bandA
		if groupOf[i]%2 == 1 {
			bg = bandB
		}

		lines := renderNode(node, expanded, width, cfg, agents)
		selected := i == cursor

		for j, line := range lines {
			if j > 0 {
				sb.WriteByte('\n')
			}
			plain := stripAnsi(line)
			if selected {
				sb.WriteString(cursorAccent)
				padded := plain + strings.Repeat(" ", max(0, width-1-lipgloss.Width(plain)))
				sb.WriteString(lipgloss.NewStyle().
					Foreground(cursorFg).
					Background(cursorBg).
					Render(padded))
			} else {
				padded := plain + strings.Repeat(" ", max(0, width-lipgloss.Width(plain)))
				sb.WriteString(applyBand(padded, node, bg))
			}
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}

func applyBand(text string, node TreeNode, bg lipgloss.Color) string {
	fg := lipgloss.Color("244")
	switch node.Kind {
	case WindowNode:
		if node.Window.Active {
			fg = lipgloss.Color("255")
		} else {
			fg = lipgloss.Color("250")
		}
	case PaneNode:
		if node.Pane.Active {
			fg = lipgloss.Color("252")
		} else {
			fg = lipgloss.Color("244")
		}
	}
	return lipgloss.NewStyle().
		Foreground(fg).
		Background(bg).
		Render(text)
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

	label := fmt.Sprintf(" %d  %s", w.Index, w.Name)

	if w.Active {
		indicator := activeIndicator
		labelWidth := lipgloss.Width(label)
		indicatorWidth := lipgloss.Width(indicator)
		pad := width - labelWidth - indicatorWidth - 1
		if pad < 1 {
			pad = 1
		}
		return []string{label + strings.Repeat(" ", pad) + "●"}
	}

	return []string{label}
}

func renderPaneNode(node TreeNode, width int, cfg config.Config, agents map[int]detect.Agent) []string {
	p := node.Pane
	var lines []string

	label := "   " + paneLabel(p)
	if _, ok := agents[p.PID]; ok {
		label += "  🤖"
	}

	lines = append(lines, label)

	if cfg.ShowCwd || cfg.ShowGit {
		var parts []string
		if cfg.ShowCwd {
			parts = append(parts, filepath.Base(p.CurrentPath))
		}
		if cfg.ShowGit {
			if branch := git.Branch(p.CurrentPath); branch != "" {
				parts = append(parts, branch)
			}
		}
		if len(parts) > 0 {
			lines = append(lines, "   "+strings.Join(parts, "  "))
		}
	}

	return lines
}

func paneLabel(p *tmux.Pane) string {
	if p.Title != "" {
		host, _ := os.Hostname()
		if p.Title != host {
			return p.Title
		}
	}
	return p.Command
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
