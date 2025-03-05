INSTALL_PATH := $(HOME)/.local/bin/td
VERSION := $(shell grep -o '"[0-9]\+\.[0-9]\+\.[0-9]\+"' main.go | tr -d '"')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(GIT_COMMIT) -X main.BuildDate=$(BUILD_DATE)"

all: lint test build

test:
	go test -v ./...

lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

build: test
	go build $(LDFLAGS) -o bin/td main.go

clean:
	rm -f bin/td

run:
	go run main.go

install: build
	cp -f bin/td $(INSTALL_PATH)

# Build for multiple platforms
build-release:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/td-$(VERSION)-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/td-$(VERSION)-darwin-amd64 main.go
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/td-$(VERSION)-windows-amd64.exe main.go

.PHONY: all build test clean run deps build-release install lint
