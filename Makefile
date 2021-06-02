MAKEFLAGS += --warn-undefined-variables
VERBOSE = -v

GOCMD = go
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

.PHONY: tools clean tidy test snapshot release

# don't `tidy` by default or `tools` will break
all: clean tools test build

$(DISTDIR):
	mkdir -p $(DISTDIR)

$(DISTDIR)/$(TARGET): $(DISTDIR) $(SOURCES) test
	$(GOBUILD) $(VERBOSE) -o $(DISTDIR)/$(TARGET) $(MAIN)

clean:
	$(GOCLEAN) $(VERBOSE)
	$(GOCLEAN) $(VERBOSE) -modcache
	rm -rf $(DISTDIR)/

tools:
	$(GOINSTALL) $(VERBOSE) honnef.co/go/tools/cmd/staticcheck
	$(GOINSTALL) $(VERBOSE) github.com/goreleaser/goreleaser

test: tools $(SOURCES)
	staticcheck ./...
	$(GOTEST) $(VERBOSE) ./...

build: $(DISTDIR)/$(TARGET)

tidy:
	$(GOMOD) tidy

release: $(SOURCES)
	goreleaser --rm-dist

snapshot: $(SOURCES)
	goreleaser --snapshot --rm-dist
