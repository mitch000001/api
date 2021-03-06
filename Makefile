COMMIT     := $(shell git rev-parse --short HEAD)
VERSION    := 0.0.1

LDFLAGS    := -ldflags \
              "-X main.Commit $(COMMIT)\
               -X main.Version $(VERSION)"

GOOS       := $(shell go env GOOS)
GOARCH     := $(shell go env GOARCH)
GOBUILD    := GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS)

ARCHIVE    := umsatz-$(VERSION)-$(GOOS)-$(GOARCH).tar.gz
DISTDIR    := dist/$(GOOS)_$(GOARCH)

.PHONY: default archive clean install

default: *.go
	$(GOBUILD)

archive: dist/$(ARCHIVE)

GO ?= go

all: build

build:
	$(GO) build

check:
	@$(GO) list -f '{{join .Deps "\n"}}' | xargs $(GO) list -f '{{if not .Standard}}{{.ImportPath}} {{.Dir}}{{end}}' | column -t

clean:
	git clean -f -x -d

dist/$(ARCHIVE): $(DISTDIR)/api
	tar -C $(DISTDIR) -czvf $@ .

$(DISTDIR)/api: *.go
	$(GOBUILD) -o $@
