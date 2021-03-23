.DEFAULT_GOAL := build

DIRS=bin
BINARY=spectrograph-worker

$(info $(shell mkdir -p $(DIRS)))
BIN=$(CURDIR)/bin
export GOBIN=$(CURDIR)/bin

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

fetch: ## Fetches the necessary dependencies to build.
	go mod download
	go mod tidy

upgrade-deps: ## Upgrade all dependencies to the latest version.
	go get -u ./...

upgrade-deps-patch: ## Upgrade all dependencies to the latest patch release.
	go get -u=patch ./...

clean: ## Cleans up generated files/folders from the build.
	/bin/rm -rfv "rice-box.go" "${BINARY}"

compile:
	go build -ldflags '-s -w' -tags netgo -installsuffix netgo -v -o "${BINARY}"

build: fetch compile ## Builds the application
	echo

debug: fetch ## Runs the application in debug mode
	go run *.go --debug \
		--log.level debug
