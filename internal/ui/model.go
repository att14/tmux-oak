package ui

import (
	"os"

	"github.com/att14/tmux-oak/internal/config"
	"github.com/att14/tmux-oak/internal/tmux"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	client    *tmux.Client
	session   string
	ownPaneID string
	cfg       config.Config
	state     *tmux.State
	expanded  map[int]bool
	nodes     []TreeNode
	cursor    int
	width     int
	height    int
	err       error
}

func NewModel(client *tmux.Client, session string, cfg config.Config) *Model {
	return &Model{
		client:    client,
		session:   session,
		ownPaneID: os.Getenv("TMUX_PANE"),
		cfg:       cfg,
		expanded:  make(map[int]bool),
	}
}

type stateMsg struct{ state *tmux.State }
type errMsg struct{ err error }
type switchedMsg struct{}

func fetchState(client *tmux.Client, session, excludeID string) tea.Cmd {
	return func() tea.Msg {
		s, err := tmux.Snapshot(client, session, excludeID)
		if err != nil {
			return errMsg{err}
		}
		return stateMsg{s}
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

func (m *Model) Init() tea.Cmd {
	return fetchState(m.client, m.session, m.ownPaneID)
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
		return m, nil

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
		}

	case key.Matches(msg, keys.Down):
		if m.cursor < len(m.nodes)-1 {
			m.cursor++
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
	return renderTree(m.nodes, m.cursor, m.expanded, w, m.cfg)
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
