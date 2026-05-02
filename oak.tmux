#!/usr/bin/env bash

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

oak_key="$(tmux show-option -gqv "@oak-key" 2>/dev/null)"
: "${oak_key:=e}"

tmux bind-key "$oak_key" run-shell "$CURRENT_DIR/scripts/toggle.sh"
