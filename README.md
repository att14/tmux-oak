# tmux-oak

A tmux sidebar plugin that displays a tree view of your windows, panes, and running AI agent sessions. Inspired by [cmux](https://cmux.com/).

## Features

- Toggleable sidebar with window/pane tree view
- Per-pane metadata: running command, working directory, git branch
- Claude Code session detection (extensible to other AI agents)
- Pane content preview on highlight
- Navigate and switch windows/panes from the sidebar

## Requirements

- tmux 3.2+
- [TPM](https://github.com/tmux-plugins/tpm)

## Install

Add to `~/.tmux.conf`:

```tmux
set -g @plugin 'att14/tmux-oak'
```

Press `prefix + I` to install.

## Usage

Press `prefix + e` to toggle the sidebar.

- **Up/Down** — navigate the tree
- **Enter** — switch to the highlighted window/pane
- **Left/Right** — collapse/expand window nodes

## Configuration

| Option | Default | Description |
|--------|---------|-------------|
| `@oak-key` | `e` | Toggle keybinding (after prefix) |
| `@oak-width` | `28` | Sidebar width in columns |
| `@oak-position` | `left` | Sidebar position: `left` or `right` |
| `@oak-show-cmd` | `on` | Show pane running command |
| `@oak-show-cwd` | `on` | Show pane working directory |
| `@oak-show-git` | `on` | Show git branch per pane |
| `@oak-show-preview` | `on` | Show pane content preview |
| `@oak-refresh` | `3` | Poll interval in seconds |

## Development

```bash
make build      # Build the binary
make install    # Symlink into TPM plugins directory
make test       # Run tests
make clean      # Remove build artifacts
```

## License

MIT
