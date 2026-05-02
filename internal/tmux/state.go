package tmux

import (
	"fmt"
	"strconv"
	"strings"
)

type Window struct {
	Index  int
	Name   string
	Active bool
	Panes  []Pane
}

type Pane struct {
	Index       int
	PID         int
	Command     string
	Title       string
	CurrentPath string
	Active      bool
	ID          string
}

type State struct {
	SessionName string
	Windows     []Window
}

func Snapshot(client *Client, session string, excludePaneID string) (*State, error) {
	out, err := client.run("list-windows", "-t", session, "-F",
		"#{window_index}\t#{window_name}\t#{window_active}")
	if err != nil {
		return nil, fmt.Errorf("list-windows: %w", err)
	}

	state := &State{SessionName: session}

	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}
		fields := strings.SplitN(line, "\t", 3)
		if len(fields) < 3 {
			continue
		}

		idx, _ := strconv.Atoi(fields[0])
		w := Window{
			Index:  idx,
			Name:   fields[1],
			Active: fields[2] == "1",
		}

		paneOut, err := client.run("list-panes",
			"-t", fmt.Sprintf("%s:%d", session, idx),
			"-F", "#{pane_index}\t#{pane_pid}\t#{pane_current_command}\t#{pane_title}\t#{pane_current_path}\t#{pane_active}\t#{pane_id}")
		if err != nil {
			continue
		}

		for _, pline := range strings.Split(paneOut, "\n") {
			if pline == "" {
				continue
			}
			pf := strings.SplitN(pline, "\t", 7)
			if len(pf) < 7 {
				continue
			}

			if excludePaneID != "" && pf[6] == excludePaneID {
				continue
			}

			if pf[2] == "oak" {
				continue
			}

			pidx, _ := strconv.Atoi(pf[0])
			ppid, _ := strconv.Atoi(pf[1])
			w.Panes = append(w.Panes, Pane{
				Index:       pidx,
				PID:         ppid,
				Command:     pf[2],
				Title:       pf[3],
				CurrentPath: pf[4],
				Active:      pf[5] == "1",
				ID:          pf[6],
			})
		}

		state.Windows = append(state.Windows, w)
	}

	return state, nil
}
