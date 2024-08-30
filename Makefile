# Makefile

BINARY_NAME := taucoder-go-client
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "no-tag")

.PHONY: all build clean

all: build

build:
	go build -o $(BINARY_NAME) -ldflags "-X main.gitCommit=$(GIT_COMMIT) -X main.gitTag=$(GIT_TAG)" main.go

clean:
	rm -f $(BINARY_NAME)

