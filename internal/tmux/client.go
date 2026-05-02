package tmux

import (
	"fmt"
	"os/exec"
	"strings"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) run(args ...string) (string, error) {
	cmd := exec.Command("tmux", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("tmux %s: %w", strings.Join(args, " "), err)
	}
	return strings.TrimRight(string(out), "\n"), nil
}

func (c *Client) SelectWindow(session string, windowIndex int) error {
	_, err := c.run("select-window", "-t", fmt.Sprintf("%s:%d", session, windowIndex))
	return err
}

func (c *Client) SelectPane(session string, windowIndex, paneIndex int) error {
	_, err := c.run("select-pane", "-t", fmt.Sprintf("%s:%d.%d", session, windowIndex, paneIndex))
	return err
}
