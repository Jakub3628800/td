INSTALL_PATH := $(HOME)/.local/bin/td
VERSION := $(shell grep -o '"[0-9]\+\.[0-9]\+\.[0-9]\+"' main.go | tr -d '"')
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE)"

all: test build

test:
	go test -v ./...

build: test
	go build $(LDFLAGS) -o bin/td main.go

build-release:
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/td-$(VERSION)-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/td-$(VERSION)-darwin-amd64 main.go
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/td-$(VERSION)-windows-amd64.exe main.go

clean:
	rm -f bin/td

run:
	go run $(LDFLAGS) main.go

install: build
	cp -f bin/td $(INSTALL_PATH)


.PHONY: all build test clean run deps build-linux build-release
