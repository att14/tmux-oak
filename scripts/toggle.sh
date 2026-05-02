#!/usr/bin/env bash

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PLUGIN_DIR="$(cd "$CURRENT_DIR/.." && pwd)"
OAK_BIN="$PLUGIN_DIR/bin/oak"

if [ ! -x "$OAK_BIN" ]; then
	tmux display-message "oak: binary not found at $OAK_BIN"
	exit 1
fi

oak_width="$(tmux show-option -gqv "@oak-width" 2>/dev/null)"
: "${oak_width:=28}"

oak_position="$(tmux show-option -gqv "@oak-position" 2>/dev/null)"
: "${oak_position:=left}"

session="$(tmux display-message -p '#{session_name}')"
current_window="$(tmux display-message -p '#{window_id}')"

oak_pane_id="$(tmux show-environment -g OAK_PANE_ID 2>/dev/null | sed 's/^OAK_PANE_ID=//')"

if [ -n "$oak_pane_id" ]; then
	pane_window="$(tmux display-message -p -t "$oak_pane_id" '#{window_id}' 2>/dev/null)"

	if [ "$pane_window" = "$current_window" ]; then
		tmux kill-pane -t "$oak_pane_id"
		tmux set-environment -gu OAK_PANE_ID
		exit 0
	fi

	tmux kill-pane -t "$oak_pane_id" 2>/dev/null
	tmux set-environment -gu OAK_PANE_ID
fi

split_flags="-h"
if [ "$oak_position" = "left" ]; then
	split_flags="-hb"
fi

pane_id="$(tmux split-window $split_flags -l "$oak_width" -P -F "#{pane_id}" \
	"$OAK_BIN --session '$session'")"

tmux set-environment -g OAK_PANE_ID "$pane_id"
