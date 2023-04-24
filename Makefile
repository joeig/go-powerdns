GOCMD=go
GOTEST=$(GOCMD) test
GOCOVER=$(GOCMD) tool cover
GOFMT=gofmt
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

.DEFAULT_GOAL := all
.PHONY: all test test-without-mocks coverage check-fmt check-fmt-list fmt

all: check-fmt test coverage

test:
	$(GOTEST) -v ./... -coverprofile=./coverage.out

test-without-mocks:
	GONOMOCKS=1 $(GOTEST) -v ./... -coverprofile=./coverage.out

coverage:
	$(GOCOVER) -func=./coverage.out

check-fmt:
	@$(GOFMT) -d ${GOFILES}

check-fmt-list:
	@$(GOFMT) -l ${GOFILES}

fmt:
	$(GOFMT) -w ${GOFILES}
