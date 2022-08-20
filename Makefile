GOCMD=go
GOTEST=$(GOCMD) test
GOCOVER=$(GOCMD) tool cover
GOFMT=gofmt
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

.DEFAULT_GOAL := all
.PHONY: all test test-without-mocks coverage check-fmt fmt

all: check-fmt test coverage

test:
	$(GOTEST) -v ./... -covermode=count -coverprofile=./coverage.out

test-without-mocks:
	GONOMOCKS=1 $(GOTEST) -v ./... -covermode=count -coverprofile=./coverage.out

coverage:
	$(GOCOVER) -func=./coverage.out

check-fmt:
	$(GOFMT) -d ${GOFILES}

fmt:
	$(GOFMT) -w ${GOFILES}
