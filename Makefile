GO=GO111MODULE=on go
GOBUILD=$(GO) build

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: vendor build check test

build:
	$(GOBUILD) ./cli/moovio_xsd2go

.PHONY: pkger vendor
pkger:
ifeq ("$(wildcard $(GOBIN)/pkger)","")
	go get -u -v github.com/markbates/pkger/cmd/pkger
endif

ci-update-bundled-deps: pkger
	$(GOBIN)/pkger -o pkg/template
	go fmt ./pkg/template

.PHONY: vendor
vendor:
	rm -rf vendor
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) mod verify

test:
	$(GO) test -v github.com/moov-io/xsd2go/... -count=1 -p 1 -parallel 1

.PHONY: check
check:
ifeq ($(OS),Windows_NT)
	@echo "Skipping checks on Windows, currently unsupported."
else
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	# COVER_THRESHOLD=60.0 ./lint-project.sh
endif
