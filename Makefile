VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
LDFLAGS  = -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)

.PHONY: build install clean test lint

build:
	go build -ldflags "$(LDFLAGS)" -o bin/oak ./cmd/oak

install: build
	@echo "Linking plugin to TPM directory..."
	@ln -sfn "$(CURDIR)" "$(HOME)/.tmux/plugins/tmux-oak"
	@echo "Run 'tmux source ~/.tmux.conf' to load the plugin"

clean:
	rm -rf bin/ dist/

test:
	go test ./...

lint:
	go vet ./...
