INSTALL_PATH := /home/jk/.local/bin/td
VERSION := $(shell grep -o '"[0-9]\+\.[0-9]\+\.[0-9]\+"' main.go | tr -d '"')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(GIT_COMMIT) -X main.BuildDate=$(BUILD_DATE)"

all: test build

test:
	go test -v ./...

build:
	go build $(LDFLAGS) -o bin/td main.go

clean:
	rm -f bin/td

run:
	go run main.go

install: build
	cp -f bin/td $(INSTALL_PATH)

.PHONY: all build test clean run install
