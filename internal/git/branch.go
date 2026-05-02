package git

import (
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var (
	cache   = make(map[string]string)
	cacheMu sync.RWMutex
)

func Branch(dir string) string {
	dir = filepath.Clean(dir)

	cacheMu.RLock()
	if b, ok := cache[dir]; ok {
		cacheMu.RUnlock()
		return b
	}
	cacheMu.RUnlock()

	cmd := exec.Command("git", "-C", dir, "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	branch := strings.TrimSpace(string(out))

	cacheMu.Lock()
	cache[dir] = branch
	cacheMu.Unlock()

	return branch
}

func ClearCache() {
	cacheMu.Lock()
	cache = make(map[string]string)
	cacheMu.Unlock()
}
