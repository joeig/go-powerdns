GOCMD=go
GOTEST=$(GOCMD) test
GOCOVER=$(GOCMD) tool cover
GOFMT=gofmt
GODEP=dep
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

.DEFAULT_GOAL := all
.PHONY: all test coverage check-fmt fmt deps

all: check-fmt test coverage

test:
	$(GOTEST) -v ./... -coverprofile=c.out

coverage:
	$(GOCOVER) -func=c.out

check-fmt:
	$(GOFMT) -d ${GOFILES}

fmt:
	$(GOFMT) -w ${GOFILES}

deps:
	$(GODEP) ensure
