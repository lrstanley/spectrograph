#!/usr/bin/make -f

.DEFAULT_GOAL := build
MAKEPID:= $(shell echo $$PPID)

DIRS=bin
BINARY=http-server

$(info $(shell mkdir -p $(DIRS)))
BIN=$(CURDIR)/bin
export GOBIN=$(CURDIR)/bin
export GOPRIVATE=github.com/lrstanley/spectrograph/*

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

fetch: ## Fetches the necessary dependencies to build.
	test -f $(BIN)/bin/rice || go get -u -v github.com/GeertJohan/go.rice/rice
	go mod download
	go mod tidy

upgrade-deps: ## Upgrade all dependencies to the latest version.
	go get -u ./...

upgrade-deps-patch: ## Upgrade all dependencies to the latest patch release.
	go get -u=patch ./...

clean: ## Cleans up generated files/folders from the build.
	/bin/rm -rfv "tmp" "public/dist/*" "rice-box.go" "${BINARY}"

generate: ## Generate go bundled files from frontend
	$(BIN)/rice -v embed-go

compile:
	go build -ldflags '-s -w' -tags netgo -installsuffix netgo -v -o "${BINARY}"

build: fetch clean generate compile ## Builds the application (with generate).
	echo

debug: ## Runs the application in debug mode
	# TODO: add -race, but add options for globally larger http timeouts (or timeouts in general?)
	go run *.go --debug \
		--http.bind-addr :8080 \
		--http.base-url "http://localhost:8081/" \
		--log.level debug \
		--mongo.db-name spectrograph-dev \
		--mongo.uri "mongodb://mongo:27017/?maxPoolSize=200" \
		--migration.purge
