#!/usr/bin/make -f

.DEFAULT_GOAL := build
BINARY = http-server
MAKEPID := $(shell echo $$PPID)


GO_MOD_ROOT = $(shell dirname $(shell go env GOMOD))
export GOBIN = ${GO_MOD_ROOT}/bin
export GOMODCACHE = ${GO_MOD_ROOT}/.gomodcache
$(info $(shell mkdir -p $(GOBIN)))
export PATH = $(shell printenv PATH):$(GOBIN)
export GOPRIVATE = github.com/lrstanley/spectrograph/*

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

fetch: ## Fetches the necessary dependencies to build.
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install github.com/twitchtv/twirp/protoc-gen-twirp
	go install github.com/favadi/protoc-go-inject-tag
	go mod download
	go mod tidy
	go generate -x github.com/lrstanley/spectrograph/...
	find . -not -path '*/\.*' -type f -name "*.pb.go" -exec protoc-go-inject-tag -input="{}" \;

upgrade-deps: ## Upgrade all dependencies to the latest version.
	go get -u ./...

upgrade-deps-patch: ## Upgrade all dependencies to the latest patch release.
	go get -u=patch ./...

clean: ## Cleans up generated files/folders from the build.
	/bin/rm -rfv "tmp" "public/dist/*" "rice-box.go" "${BINARY}"

generate: ## Generate go bundled files from frontend
	$(GOBIN)/rice -v embed-go

compile:
	go build -mod=readonly -ldflags '-s -w' -tags netgo -installsuffix netgo -v -o "${BINARY}"

build: fetch clean generate compile ## Builds the application (with generate).
	echo

debug: fetch ## Runs the application in debug mode
	# TODO: add -race, but add options for globally larger http timeouts (or timeouts in general?)
	go run *.go --debug \
		--http.bind-addr :8080 \
		--http.base-url "http://localhost:8081/" \
		--log.level debug \
		--mongo.db-name spectrograph-dev \
		--mongo.uri "mongodb://mongo:27017/?maxPoolSize=200" \
		--migration.purge
