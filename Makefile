BINARY := keep
VERSION ?= dev
PREFIX ?= $(HOME)/.local
BINDIR ?= $(PREFIX)/bin

.PHONY: build install install-wsl test clean

build:
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o bin/$(BINARY) ./cmd/keep

install: build
	mkdir -p "$(BINDIR)"
	cp "bin/$(BINARY)" "$(BINDIR)/$(BINARY)"
	@echo "Installed $(BINARY) to $(BINDIR)/$(BINARY)"

install-wsl:
	@mkdir -p "$(BINDIR)"
	@case "$$(uname -m)" in \
		x86_64) goarch=amd64 ;; \
		aarch64|arm64) goarch=arm64 ;; \
		*) echo "Unsupported WSL architecture: $$(uname -m)" >&2; exit 1 ;; \
	esac; \
	GOOS=windows GOARCH=$$goarch go build -ldflags "-s -w -X main.version=$(VERSION)" -o "$(BINDIR)/keep.exe" ./cmd/keep
	@printf '%s\n' '#!/bin/sh' 'exec "$$(dirname "$$0")/keep.exe" "$$@"' > "$(BINDIR)/keep"
	@chmod +x "$(BINDIR)/keep"
	@echo "Installed Windows keep.exe and WSL wrapper to $(BINDIR)"

test:
	go test ./...

clean:
	rm -rf bin
