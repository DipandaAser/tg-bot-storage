.DEFAULT_GOAL := help

## test: run tests on cmd and pkg files.
.PHONY: test
test: vet fmt
	go test ./...

## fmt: format the code
.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

check-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2
	go install github.com/kisielk/errcheck@latest

## lint: run linters over the entire code base
.PHONY: lint
lint: check-lint
	golangci-lint run ./... --timeout 15m0s
	errcheck -exclude ./.golangci-errcheck-exclude.txt ./...

all: help
.PHONY: help
help: Makefile
	@echo " Choose a command..."
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'