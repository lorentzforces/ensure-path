SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
.SILENT:

help:
	echo '  make clean        remove generated files'
	echo '  make build        build the project from scratch'
	echo '  make ensure-path  build executable if not already built'
	echo '  make test         execute tests and checks'
.PHONY: help

# go builds are fast enough that we can just build on demand instead of trying to do any fancy
# change detection
build: clean ensure-path
.PHONY: build

ensure-path:
	go build ./cmd/ensure-path

clean:
	rm -f ./ensure-path
.PHONY: clean

test:
	go test ./...
.PHONY: test
