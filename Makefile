BINARY := stay
VERSION ?= dev
PREFIX ?= $(HOME)/.local
BINDIR ?= $(PREFIX)/bin

.PHONY: build install test clean

build:
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o bin/$(BINARY) ./cmd/stay

install: build
	mkdir -p "$(BINDIR)"
	cp "bin/$(BINARY)" "$(BINDIR)/$(BINARY)"
	@echo "Installed $(BINARY) to $(BINDIR)/$(BINARY)"

test:
	go test ./...

clean:
	rm -rf bin
