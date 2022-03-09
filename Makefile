#@IgnoreInspection BashAddShebang

export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

export APP=gophermq

export LDFLAGS="-w -s"

all: format lint build

run-broker:
	go run -ldflags $(LDFLAGS) ./cmd/$(APP) broker

build:
	go build -ldflags $(LDFLAGS) ./cmd/$(APP)

install:
	go install -ldflags $(LDFLAGS) ./cmd/$(APP)

check-formatter:
	which goimports || GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

format: check-formatter
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R goimports -w R
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R gofmt -s -w R

check-linter:
	which golangci-lint || GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint@v1.23.8

lint: check-linter
	golangci-lint run $(ROOT)/...

test:
	go test -ldflags $(LDFLAGS) -v -race -p 1 `go list ./... | grep -v integration`

ci-test:
	go test -ldflags $(LDFLAGS) -v -race -p 1 -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -func coverage.txt

integration-tests:
	go test -ldflags $(LDFLAGS) -v -race -p 1 `go list ./... | grep integration`