GOCMD=go
GOTEST=$(GOCMD) test
GOFMT=gofmt
GODEP=dep
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

.DEFAULT_GOAL := all
.PHONY: all test check-fmt fmt deps

all: check-fmt fmt test

test:
	$(GOTEST) -v ./...

check-fmt:
	$(GOFMT) -d ${GOFILES}

fmt:
	$(GOFMT) -w ${GOFILES}

deps:
	$(GODEP) ensure
