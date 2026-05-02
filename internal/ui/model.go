package ui

import (
	"fmt"
	"os"
	"time"

	"github.com/att14/tmux-oak/internal/config"
	"github.com/att14/tmux-oak/internal/detect"
	"github.com/att14/tmux-oak/internal/git"
	"github.com/att14/tmux-oak/internal/tmux"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	client    *tmux.Client
	session   string
	ownPaneID string
	cfg       config.Config
	registry  *detect.Registry
	state     *tmux.State
	agents      map[int]detect.Agent
	previewText string
	previewTgt  string
	expanded    map[int]bool
	nodes     []TreeNode
	cursor    int
	width     int
	height    int
	err       error
}

func NewModel(client *tmux.Client, session string, cfg config.Config, registry *detect.Registry) *Model {
	return &Model{
		client:    client,
		session:   session,
		ownPaneID: os.Getenv("TMUX_PANE"),
		cfg:       cfg,
		registry:  registry,
		expanded:  make(map[int]bool),
		agents:    make(map[int]detect.Agent),
	}
}

type stateMsg struct{ state *tmux.State }
type errMsg struct{ err error }
type switchedMsg struct{}
type tickMsg time.Time
type previewMsg struct {
	content string
	target  string
}

func fetchState(client *tmux.Client, session, excludeID string) tea.Cmd {
	return func() tea.Msg {
		s, err := tmux.Snapshot(client, session, excludeID)
		if err != nil {
			return errMsg{err}
		}
		return stateMsg{s}
	}
}

func capturePreview(client *tmux.Client, session string, node TreeNode) tea.Cmd {
	if node.Kind != PaneNode {
		return nil
	}
	wIdx := node.WindowIndex
	pIdx := node.PaneIndex
	target := fmt.Sprintf("%d.%d", wIdx, pIdx)
	return func() tea.Msg {
		content, err := client.CapturePane(session, wIdx, pIdx, 20)
		if err != nil {
			return previewMsg{content: "", target: target}
		}
		return previewMsg{content: content, target: target}
	}
}

func doSwitch(client *tmux.Client, session string, node TreeNode) tea.Cmd {
	wIdx := node.WindowIndex
	pIdx := node.PaneIndex
	isPane := node.Kind == PaneNode
	return func() tea.Msg {
		_ = client.SelectWindow(session, wIdx)
		if isPane {
			_ = client.SelectPane(session, wIdx, pIdx)
		}
		return switchedMsg{}
	}
}

func tickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		fetchState(m.client, m.session, m.ownPaneID),
		tickCmd(time.Duration(m.cfg.Refresh)*time.Second),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case stateMsg:
		firstLoad := m.state == nil
		m.state = msg.state

		panePIDs := make(map[int]bool)
		for _, w := range m.state.Windows {
			for _, p := range w.Panes {
				panePIDs[p.PID] = true
			}
		}
		m.agents = m.registry.Scan(panePIDs)

		if firstLoad {
			for _, w := range m.state.Windows {
				if w.Active {
					m.expanded[w.Index] = true
				}
			}
		}
		m.rebuildNodes()
		if firstLoad {
			for i, n := range m.nodes {
				if n.Kind == WindowNode && n.Window.Active {
					m.cursor = i
					break
				}
			}
		}
		return m, m.previewCurrent()

	case previewMsg:
		m.previewText = msg.content
		m.previewTgt = msg.target
		return m, nil

	case tickMsg:
		git.ClearCache()
		return m, tea.Batch(
			fetchState(m.client, m.session, m.ownPaneID),
			tickCmd(time.Duration(m.cfg.Refresh)*time.Second),
		)

	case switchedMsg:
		return m, fetchState(m.client, m.session, m.ownPaneID)

	case errMsg:
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, keys.Up):
		if m.cursor > 0 {
			m.cursor--
			return m, m.previewCurrent()
		}

	case key.Matches(msg, keys.Down):
		if m.cursor < len(m.nodes)-1 {
			m.cursor++
			return m, m.previewCurrent()
		}

	case key.Matches(msg, keys.Enter):
		if m.cursor >= 0 && m.cursor < len(m.nodes) {
			return m, doSwitch(m.client, m.session, m.nodes[m.cursor])
		}

	case key.Matches(msg, keys.Right):
		if m.cursor >= 0 && m.cursor < len(m.nodes) {
			node := m.nodes[m.cursor]
			if node.Kind == WindowNode && !m.expanded[node.WindowIndex] {
				m.expanded[node.WindowIndex] = true
				m.rebuildNodes()
			}
		}

	case key.Matches(msg, keys.Left):
		if m.cursor >= 0 && m.cursor < len(m.nodes) {
			node := m.nodes[m.cursor]
			winIdx := node.WindowIndex
			if m.expanded[winIdx] {
				delete(m.expanded, winIdx)
				m.rebuildNodes()
				for i, n := range m.nodes {
					if n.Kind == WindowNode && n.WindowIndex == winIdx {
						m.cursor = i
						break
					}
				}
			}
		}
	}

	return m, nil
}

func (m *Model) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n\nPress q to quit."
	}
	if m.state == nil {
		return "Loading..."
	}
	w := m.width
	if w == 0 {
		w = 28
	}
	view := renderTree(m.nodes, m.cursor, m.expanded, w, m.cfg, m.agents)
	if m.cfg.ShowPreview && m.previewText != "" {
		view += "\n" + renderPreview(m.previewText, m.previewTgt, w, previewLines)
	}
	return view
}

func (m *Model) previewCurrent() tea.Cmd {
	if !m.cfg.ShowPreview {
		return nil
	}
	if m.cursor >= 0 && m.cursor < len(m.nodes) {
		node := m.nodes[m.cursor]
		if node.Kind == PaneNode {
			return capturePreview(m.client, m.session, node)
		}
	}
	m.previewText = ""
	m.previewTgt = ""
	return nil
}

func (m *Model) rebuildNodes() {
	m.nodes = buildNodes(m.state, m.expanded)
	if m.cursor >= len(m.nodes) {
		m.cursor = len(m.nodes) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}
