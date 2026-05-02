package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/att14/tmux-oak/internal/config"
	"github.com/att14/tmux-oak/internal/detect"
	"github.com/att14/tmux-oak/internal/tmux"
	"github.com/att14/tmux-oak/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	session := flag.String("session", "", "tmux session name")
	focusPane := flag.String("focus-pane", "", "pane ID to highlight initially")
	flag.Parse()

	if *showVersion {
		fmt.Printf("oak %s (%s)\n", version, commit)
		os.Exit(0)
	}

	if *session == "" {
		fmt.Fprintln(os.Stderr, "oak: --session is required")
		os.Exit(1)
	}

	cfg := config.Load()
	client := tmux.NewClient()

	registry := detect.NewRegistry()
	registry.Register(detect.NewClaudeDetector())

	model := ui.NewModel(client, *session, cfg, registry, *focusPane)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "oak: %v\n", err)
		os.Exit(1)
	}
}
