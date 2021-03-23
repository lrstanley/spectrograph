#!/usr/bin/make -f

.DEFAULT_GOAL := build
MAKEPID:= $(shell echo $$PPID)

DIRS=bin
BINARY=spectrograph-http

$(info $(shell mkdir -p $(DIRS)))
BIN=$(CURDIR)/bin
export GOBIN=$(CURDIR)/bin

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

fetch-go: ## Fetches the necessary dependencies to build.
	test -f $(BIN)/bin/rice || go get -u -v github.com/GeertJohan/go.rice/rice
	go mod download
	go mod tidy

fetch-node: ## Fetches the necessary NodeJS dependencies to build.
	test -d public/node_modules || (cd public && npm install)

upgrade-deps: ## Upgrade all dependencies to the latest version.
	go get -u ./...

upgrade-deps-patch: ## Upgrade all dependencies to the latest patch release.
	go get -u=patch ./...

clean: ## Cleans up generated files/folders from the build.
	/bin/rm -rfv "tmp" "public/dist" "rice-box.go" "${BINARY}"

clean-cache: ## Cleans up generated cache (speeds up during dev time).
	/bin/rm -rfv "public/.cache"

generate-watch: ## Generate public html/css/js when files change (faster, but larger files).
	cd public && npm run watch

generate-node: ## Generate public html/css/js files for use in production (slower, smaller/minified files).
	cd public && npm run build

generate-go: ## Generate go bundled files from frontend
	$(BIN)/rice -v embed-go

compile:
	go build -ldflags '-s -w' -tags netgo -installsuffix netgo -v -o "${BINARY}"

build: fetch-go fetch-node clean clean-cache generate-node generate-go compile ## Builds the application (with generate).
	echo

debug: ## Runs the application in debug mode
	# TODO: add -race, but add options for globally larger http timeouts (or timeouts in general?)
	go run *.go --debug --http :8080 \
		--log.level debug \
		--mongo.db-name spectrograph-dev \
		--mongo.uri "mongodb://mongo:27017/?maxPoolSize=200" \
		--migration.purge
