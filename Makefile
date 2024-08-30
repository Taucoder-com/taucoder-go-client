# Makefile

BINARY_NAME := taucoder-go-client
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

.PHONY: all build clean

all: build

build:
	go build -o $(BINARY_NAME) -ldflags "-X main.buildTimestamp=$(BUILD_DATE)" main.go

clean:
	rm -f $(BINARY_NAME)

