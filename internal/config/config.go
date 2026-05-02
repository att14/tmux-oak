package config

import (
	"os/exec"
	"strings"
)

type Config struct {
	ShowCmd     bool
	ShowCwd     bool
	ShowGit     bool
	ShowPreview bool
	Width       int
	Position    string
	Refresh     int
}

func Load() Config {
	return Config{
		ShowCmd:     optionBool("@oak-show-cmd", true),
		ShowCwd:     optionBool("@oak-show-cwd", true),
		ShowGit:     optionBool("@oak-show-git", true),
		ShowPreview: optionBool("@oak-show-preview", false),
		Width:       28,
		Position:    optionStr("@oak-position", "left"),
		Refresh:     3,
	}
}

func optionStr(name, fallback string) string {
	cmd := exec.Command("tmux", "show-option", "-gqv", name)
	out, err := cmd.Output()
	if err != nil {
		return fallback
	}
	val := strings.TrimSpace(string(out))
	if val == "" {
		return fallback
	}
	return val
}

func optionBool(name string, fallback bool) bool {
	val := optionStr(name, "")
	if val == "" {
		return fallback
	}
	return val == "on" || val == "true" || val == "1"
}
