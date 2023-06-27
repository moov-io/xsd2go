USERID := $(shell id -u $$USER)
GROUPID:= $(shell id -g $$USER)

all: vendor build check test

build:
	go install golang.org/x/tools/cmd/goimports@latest
	go build ./cli/gocomply_xsd2go

.PHONY: vendor
vendor:
	rm -rf vendor
	go mod tidy
	go mod vendor
	go mod verify

test:
	go test -v github.com/gocomply/xsd2go/... -count=1 -p 1 -parallel 1

.PHONY: check
check:
ifeq ($(OS),Windows_NT)
	@echo "Skipping checks on Windows, currently unsupported."
else
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	# COVER_THRESHOLD=60.0 ./lint-project.sh
endif
