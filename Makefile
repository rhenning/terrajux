GOCMD = go
GOMOD = $(GOCMD) mod
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
CMDDIR = cmd
DISTDIR = dist
TARGET = terrajux
MAIN = $(CMDDIR)/$(TARGET)/main.go
SOURCES = $(wildcard **/*.go)

.PHONY: test clean tidy release snapshot

$(DISTDIR)/$(TARGET): test $(SOURCES)
	$(GOBUILD) -o $(DISTDIR)/$(TARGET) -v $(MAIN)

test: tidy $(SOURCES)
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	mkdir -p $(DISTDIR)
	rm -f $(DISTDIR)/$(TARGET)

tidy:
	$(GOMOD) tidy

release: $(SOURCES)
	goreleaser --rm-dist

snapshot: $(SOURCES)
	goreleaser --snapshot --rm-dist
