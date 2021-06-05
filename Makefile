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

.PHONY: tools clean tidy test unit static snapshot release

# don't `tidy` by default or `tools` will break
all: clean tools test build

$(DISTDIR):
	mkdir -p $(DISTDIR)

$(DISTDIR)/$(TARGET): $(DISTDIR) $(SOURCES) test
	$(GOBUILD) $(VERBOSE) -o $(DISTDIR)/$(TARGET) $(MAIN)

clean:
	$(GOCLEAN) $(VERBOSE)
	$(GOCLEAN) $(VERBOSE) -testcache
	# please don't remove the ./ prefix here. run the build from inside this
	# directory and avoid inadvertently nuking / due to a bug.
	rm -rf ./$(DISTDIR)/

tools: $(DISTDIR)/tools-stamp

$(DISTDIR)/tools-stamp:
	$(GOINSTALL) $(VERBOSE) honnef.co/go/tools/cmd/staticcheck
	$(GOINSTALL) $(VERBOSE) github.com/securego/gosec/v2/cmd/gosec
	$(GOINSTALL) $(VERBOSE) github.com/goreleaser/goreleaser
	touch $@

test: tools static unit

static: tools $(SOURCES)
ifneq ($(strip $(GOFMTL)),)
	$(error invalid formatting detected. please run `go fmt ./...`)
endif
	staticcheck ./...
	gosec ./...

unit: tools $(SOURCES)
	$(GOTEST) $(VERBOSE) ./...

build: $(DISTDIR)/$(TARGET)

tidy:
	$(GOMOD) tidy

release: $(SOURCES)
	goreleaser --rm-dist

snapshot: $(SOURCES)
	goreleaser --snapshot --rm-dist
