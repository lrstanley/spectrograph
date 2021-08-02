.DEFAULT_GOAL := build
BINARY = spectrograph-worker

export GOBIN=$(shell dirname $(shell go env GOMOD))/bin
$(info $(shell mkdir -p $(GOBIN)))
export PATH=$(shell printenv PATH):$(GOBIN)
export GOPRIVATE=github.com/lrstanley/spectrograph/*

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

fetch: ## Fetches the necessary dependencies to build.
	go install github.com/twitchtv/twirp/protoc-gen-twirp
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go mod download
	go mod tidy
	go generate github.com/lrstanley/spectrograph/...

upgrade-deps: ## Upgrade all dependencies to the latest version.
	go get -u ./...

upgrade-deps-patch: ## Upgrade all dependencies to the latest patch release.
	go get -u=patch ./...

clean: ## Cleans up generated files/folders from the build.
	/bin/rm -rfv "rice-box.go" "${BINARY}"

compile:
	go build -mod=readonly -ldflags '-s -w' -tags netgo -installsuffix netgo -v -o "${BINARY}"

build: fetch compile ## Builds the application
	echo

debug: fetch ## Runs the application in debug mode
	go run *.go --debug \
		--rpc.uri "http://http-server:8080" \
		--discord.shard-id 0 \
		--discord.num-shards 1
