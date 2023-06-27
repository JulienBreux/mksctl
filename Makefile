NAME    ?= mksctl
BINDIR  := $(CURDIR)/bin
VERSION ?= dev
COMMIT  ?= $(shell git rev-parse HEAD)
DATE    ?= $(shell date +%FT%T%Z)

lint: ## [project] Lint
	golangci-lint run

test: ## [project] Test packages
	go test -count=1 -cover -coverprofile=cover.out -v ./...

test-cover: ## [project] Display coverage
	go tool cover -html cover.out -o cover.html

upgrade-deps: ## [project] Upgrade all dependencies
	go get -u ./...

build: ## [binary] Build local binary
	mkdir -p ./bin
	CGO_ENABLED=0 go build -ldflags='-s -w -X github.com/JulienBreux/${NAME}/pkg/version.Version=${VERSION} -X github.com/JulienBreux/${NAME}/pkg/version.Commit=${COMMIT} -X github.com/JulienBreux/${NAME}/pkg/version.RawDate=${DATE}' -o ${BINDIR}/app ./cmd/${NAME}

run: build ## [binary] Run local binary
	./bin/app

image-build: ## [image] Build local image
	docker build --build-arg VERSION=${VERSION} --build-arg COMMIT=${COMMIT} --build-arg DATE=${DATE} --no-cache -t ghcr.io/julienbreux/${NAME}:latest .

image-run: ## [image] Run local image
	docker run -p 8080:8080 -t ghcr.io/julienbreux/${NAME}:latest

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: image-build
