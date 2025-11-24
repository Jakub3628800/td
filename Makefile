INSTALL_PATH := /home/jk/.local/bin/td
VERSION := $(shell grep -o '"[0-9]\+\.[0-9]\+\.[0-9]\+"' main.go | tr -d '"')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(GIT_COMMIT) -X main.BuildDate=$(BUILD_DATE)"

all: sqlc test build

sqlc:
	sqlc generate --file sqlc/sqlc.yaml

test:
	go test -v ./...

build:
	go build $(LDFLAGS) -o bin/td main.go

build-release:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/td-$(VERSION)-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/td-$(VERSION)-darwin-amd64 main.go
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/td-$(VERSION)-windows-amd64.exe main.go

clean:
	rm -f bin/td

run:
	go run main.go

install: build
	cp -f bin/td $(INSTALL_PATH)

bump-version:
	@if [ -z "$(NEW_VERSION)" ]; then \
		echo "Usage: make bump-version NEW_VERSION=X.Y.Z"; \
		exit 1; \
	fi
	@if ! echo "$(NEW_VERSION)" | grep -qE '^[0-9]+\.[0-9]+\.[0-9]+$$'; then \
		echo "Error: Version must be in format X.Y.Z (e.g., 1.2.3)"; \
		exit 1; \
	fi
	@sed -i 's/Version = "[0-9]*\.[0-9]*\.[0-9]*"/Version = "$(NEW_VERSION)"/' main.go
	@echo "Bumped version to $(NEW_VERSION)"
	@echo "Don't forget to create a git tag: git tag v$(NEW_VERSION)"

.PHONY: all build build-release test clean run install bump-version
