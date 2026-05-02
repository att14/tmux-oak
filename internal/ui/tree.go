package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/att14/tmux-oak/internal/config"
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

func renderTree(nodes []TreeNode, cursor int, expanded map[int]bool, width int, cfg config.Config) string {
	var sb strings.Builder

	sb.WriteString(headerStyle.Render("oak"))
	sb.WriteByte('\n')
	sb.WriteString(separatorStyle.Render(strings.Repeat("━", width)))
	sb.WriteByte('\n')

	for i, node := range nodes {
		line := renderNode(node, expanded, width, cfg)
		if i == cursor {
			// Re-render with selection highlight on each line
			for j, l := range strings.Split(line, "\n") {
				if j > 0 {
					sb.WriteByte('\n')
				}
				sb.WriteString(selectedStyle.Width(width).Render(stripAnsi(l)))
			}
		} else {
			sb.WriteString(line)
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}

func renderNode(node TreeNode, expanded map[int]bool, width int, cfg config.Config) string {
	switch node.Kind {
	case WindowNode:
		return renderWindowNode(node, expanded, width)
	case PaneNode:
		return renderPaneNode(node, width, cfg)
	}
	return ""
}

func renderWindowNode(node TreeNode, expanded map[int]bool, width int) string {
	w := node.Window
	arrow := "▸"
	if expanded[w.Index] {
		arrow = "▼"
	}

	label := fmt.Sprintf(" %s %d:%s", arrow, w.Index, w.Name)

	if w.Active {
		indicator := lipgloss.NewStyle().Foreground(activeColor).Render("●")
		labelWidth := lipgloss.Width(label)
		indicatorWidth := lipgloss.Width(indicator)
		pad := width - labelWidth - indicatorWidth
		if pad < 1 {
			pad = 1
		}
		return windowStyle.Render(label) + strings.Repeat(" ", pad) + indicator
	}

	return windowDimStyle.Render(label)
}

func renderPaneNode(node TreeNode, width int, cfg config.Config) string {
	p := node.Pane
	connector := "├"
	if node.IsLastChild {
		connector = "└"
	}

	var lines []string

	// Main pane line: command name
	cmdLabel := p.Command
	if cfg.ShowCmd {
		cmdLabel = p.Command
	}
	mainLine := fmt.Sprintf("   %s %s", connector, cmdLabel)
	if p.Active {
		lines = append(lines, paneActiveStyle.Render(mainLine))
	} else {
		lines = append(lines, paneDimStyle.Render(mainLine))
	}

	// Metadata lines (indented under the pane)
	hasMetadata := cfg.ShowCwd || cfg.ShowGit
	if hasMetadata {
		continuation := "│"
		if node.IsLastChild {
			continuation = " "
		}

		var meta []string
		if cfg.ShowCwd {
			meta = append(meta, shortenPath(p.CurrentPath))
		}
		if cfg.ShowGit {
			if branch := git.Branch(p.CurrentPath); branch != "" {
				meta = append(meta, metaBranchStyle.Render(branch))
			}
		}

		if len(meta) > 0 {
			metaLine := fmt.Sprintf("   %s  %s", continuation, strings.Join(meta, "  "))
			lines = append(lines, metaDimStyle.Render(metaLine))
		}
	}

	return strings.Join(lines, "\n")
}

func shortenPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Base(path)
	}
	if strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}
	return path
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
