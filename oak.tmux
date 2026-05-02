#!/usr/bin/env bash

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OAK_BIN="$CURRENT_DIR/bin/oak"
VERSION_FILE="$CURRENT_DIR/.oak-version"
REPO="att14/tmux-oak"

install_binary() {
	local os arch
	os="$(uname -s | tr '[:upper:]' '[:lower:]')"
	arch="$(uname -m)"

	case "$arch" in
		x86_64) arch="amd64" ;;
		aarch64|arm64) arch="arm64" ;;
		*) tmux display-message "oak: unsupported architecture: $arch"; return 1 ;;
	esac

	local latest
	latest="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" 2>/dev/null | grep '"tag_name"' | sed 's/.*"v\(.*\)".*/\1/')"

	if [ -z "$latest" ]; then
		tmux display-message "oak: could not determine latest version"
		return 1
	fi

	if [ -x "$OAK_BIN" ] && [ -f "$VERSION_FILE" ]; then
		local installed
		installed="$(cat "$VERSION_FILE")"
		if [ "$installed" = "$latest" ]; then
			return 0
		fi
	fi

	local url="https://github.com/$REPO/releases/download/v${latest}/tmux-oak_${latest}_${os}_${arch}.tar.gz"

	tmux display-message "oak: downloading v${latest}..."

	mkdir -p "$CURRENT_DIR/bin"
	if curl -fsSL "$url" | tar xz -C "$CURRENT_DIR/bin" oak 2>/dev/null; then
		chmod +x "$OAK_BIN"
		echo "$latest" > "$VERSION_FILE"
		tmux display-message "oak: installed v${latest}"
	else
		tmux display-message "oak: download failed — run 'make build' to build from source"
		return 1
	fi
}

if [ ! -x "$OAK_BIN" ]; then
	install_binary
fi

oak_key="$(tmux show-option -gqv "@oak-key" 2>/dev/null)"
: "${oak_key:=e}"

tmux bind-key "$oak_key" run-shell "$CURRENT_DIR/scripts/toggle.sh"
