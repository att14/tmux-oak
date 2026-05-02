package detect

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

type ClaudeDetector struct {
	sessionsDir string
}

func NewClaudeDetector() *ClaudeDetector {
	home, _ := os.UserHomeDir()
	return &ClaudeDetector{
		sessionsDir: filepath.Join(home, ".claude", "sessions"),
	}
}

func (d *ClaudeDetector) Name() string { return "claude" }
func (d *ClaudeDetector) Icon() string { return "🤖" }

type claudeSession struct {
	PID       int    `json:"pid"`
	SessionID string `json:"sessionId"`
	Kind      string `json:"kind"`
}

func (d *ClaudeDetector) Scan(panePIDs map[int]bool) (map[int]Agent, error) {
	entries, err := os.ReadDir(d.sessionsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read sessions dir: %w", err)
	}

	result := make(map[int]Agent)

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(d.sessionsDir, entry.Name()))
		if err != nil {
			continue
		}

		var sess claudeSession
		if err := json.Unmarshal(data, &sess); err != nil {
			continue
		}

		if sess.Kind != "interactive" || sess.PID == 0 {
			continue
		}

		if !isRunning(sess.PID) {
			continue
		}

		// Walk up the process tree (up to 5 levels) to find the pane PID.
		// Claude may be several levels deep: pane shell → node → claude,
		// or pane shell → zsh → wrapper → claude, etc.
		pid := sess.PID
		for range 5 {
			ppid := getppid(pid)
			if ppid <= 1 {
				break
			}
			if panePIDs[ppid] {
				result[ppid] = Agent{Name: d.Name(), Icon: d.Icon()}
				break
			}
			pid = ppid
		}
	}

	return result, nil
}

func isRunning(pid int) bool {
	return syscall.Kill(pid, 0) == nil
}

func getppid(pid int) int {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err == nil {
		return parseProcStat(data)
	}
	out, err := exec.Command("ps", "-o", "ppid=", "-p", strconv.Itoa(pid)).Output()
	if err != nil {
		return 0
	}
	ppid, _ := strconv.Atoi(strings.TrimSpace(string(out)))
	return ppid
}

func parseProcStat(data []byte) int {
	s := string(data)
	idx := strings.LastIndex(s, ") ")
	if idx < 0 {
		return 0
	}
	fields := strings.Fields(s[idx+2:])
	if len(fields) < 2 {
		return 0
	}
	ppid, _ := strconv.Atoi(fields[1])
	return ppid
}
