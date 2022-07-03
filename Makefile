APPNAME=tg-bot-storage
.DEFAULT_GOAL := help

## test: run tests on cmd and pkg files.
.PHONY: test
test: vet fmt
	go test ./...

## build: build application binary.
.PHONY: build
build:
	go build -o $(APPNAME)

## docker-build: build the api docker image
.PHONY: docker-build
docker-build:
	docker build -t bukela/bot-storage .

## docker-run: run the api docker container
.PHONY: docker-run
docker-run:
	docker run -p 7000:7000 --env-file .env bukela/bot-storage

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