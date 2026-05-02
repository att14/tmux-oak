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

enabled="$(tmux show-environment -g OAK_ENABLED 2>/dev/null | sed 's/^OAK_ENABLED=//')"

if [ "$enabled" = "1" ]; then
	# Toggle OFF: kill all oak sidebar panes in this session
	while IFS= read -r line; do
		pane_id="${line%% *}"
		marker="${line##* }"
		if [ "$marker" = "on" ] && [ -n "$pane_id" ]; then
			tmux kill-pane -t "$pane_id" 2>/dev/null
		fi
	done < <(tmux list-panes -s -t "$session" -F '#{pane_id} #{@oak-sidebar}')
	tmux set-environment -gu OAK_ENABLED
else
	# Toggle ON: create oak sidebar in every window
	# -f = full-width split (spans entire viewport height)
	# -b = before (left side)
	# -h = horizontal (vertical split)
	split_flags="-hf"
	if [ "$oak_position" = "left" ]; then
		split_flags="-hbf"
	fi

	current_win="$(tmux display-message -p '#{window_index}')"
	prev_pane="$(tmux display-message -p '#{pane_id}')"

	while IFS= read -r win; do
		d_flag="-d"
		focus_flag=""
		if [ "$win" = "$current_win" ]; then
			d_flag=""
			focus_flag="--focus-pane '$prev_pane'"
		fi

		pane_id="$(tmux split-window $split_flags $d_flag -t "$session:$win" -l "$oak_width" -P -F "#{pane_id}" \
			"$OAK_BIN --session '$session' $focus_flag")"
		tmux set-option -p -t "$pane_id" @oak-sidebar on 2>/dev/null
	done < <(tmux list-windows -t "$session" -F '#{window_index}')

	tmux set-environment -g OAK_ENABLED 1
fi
