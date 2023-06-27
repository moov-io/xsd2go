GO=GO111MODULE=on go
GOBUILD=$(GO) build

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: vendor build test

build:
	go install golang.org/x/tools/cmd/goimports@latest
	$(GOBUILD) ./cli/gocomply_xsd2go

vendor:
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) mod verify

test:
	go test -v github.com/gocomply/xsd2go/... -count=1 -p 1 -parallel 1
