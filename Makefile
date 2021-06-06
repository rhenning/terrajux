# `make`, with no args, should do most of the things, assuming go is installed.
# it must be run from the project repository's base directory.

MAKEFLAGS += --warn-undefined-variables
VERBOSE = -v

GOCMD = go
GOFMT = gofmt
GOMOD = $(GOCMD) mod
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOINSTALL = $(GOCMD) install
CMDDIR = cmd
DISTDIR = dist
TARGET = terrajux
MAIN = $(CMDDIR)/$(TARGET)/main.go
SOURCES = $(wildcard **/*.go)
GOFMTL = $(shell $(GOFMT) -l .)

.PHONY: tools clean tidy test unit static snapshot release check

# don't `tidy` by default or `tools` will break
all: check clean tools test build

check:
ifeq ($(strip $(DISTDIR)),)
	$(error DISTDIR must be a non-empty value. a bug may exist in the Makefile)
endif
	$(GOCMD) version

$(DISTDIR):
	mkdir -p $(DISTDIR)

$(DISTDIR)/$(TARGET): $(DISTDIR) $(SOURCES) test
	$(GOBUILD) $(VERBOSE) -o $(DISTDIR)/$(TARGET) $(MAIN)

clean: check
	$(GOCLEAN) $(VERBOSE)
	$(GOCLEAN) $(VERBOSE) -testcache
	rm -rf ./$(DISTDIR)/

tools: check $(DISTDIR)/tools-stamp

$(DISTDIR)/tools-stamp: $(DISTDIR)
	$(GOINSTALL) $(VERBOSE) github.com/securego/gosec/v2/cmd/gosec
ifneq ($(strip $(CI)),true)
	$(GOINSTALL) $(VERBOSE) github.com/goreleaser/goreleaser
	$(GOINSTALL) $(VERBOSE) honnef.co/go/tools/cmd/staticcheck
else
	$(warning CI=$(CI): skip installing staticcheck+goreleaser, will run as GitHub Action)
endif
	touch $@

test: tools static unit

static: check tools $(SOURCES)
ifneq ($(strip $(GOFMTL)),)
	$(error invalid formatting detected. please run `go fmt ./...`)
endif
ifneq ($(strip $(CI)),true)
	golangci-lint run || staticcheck ./...
else
	$(warning CI=$(CI): skip golangci-lint+staticcheck, will run as GitHub Action)
endif
	gosec ./...

unit: check tools $(SOURCES)
	$(GOTEST) $(VERBOSE) ./...

build: check $(DISTDIR)/$(TARGET)

tidy: check
	$(GOMOD) tidy

release: tools $(SOURCES)
	goreleaser --rm-dist

snapshot: tools $(SOURCES)
	goreleaser --snapshot --rm-dist
