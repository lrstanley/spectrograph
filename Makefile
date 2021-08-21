#!/usr/bin/make -f

# .ONESHELL:
# .SHELLFLAGS = -e

MAKEPID := $(shell echo $$PPID)


GO_MOD_ROOT = $(shell dirname $(shell go env GOMOD))
export GOBIN = ${GO_MOD_ROOT}/bin
export GOMODCACHE = ${GO_MOD_ROOT}/.gomodcache
$(info $(shell mkdir -p $(GOBIN)))
export PATH = $(shell printenv PATH):$(GOBIN)
export GOPRIVATE = github.com/lrstanley/spectrograph/*

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

fetch-go: ## Fetches the necessary dependencies to build.
	go install github.com/GeertJohan/go.rice/rice
	go mod download
	go mod tidy
	go generate -x github.com/lrstanley/spectrograph/...

upgrade-deps-go: ## Upgrade all dependencies to the latest version.
	go get -u ./...

upgrade-deps-go-patch: ## Upgrade all dependencies to the latest patch release.
	go get -u=patch ./...

watch-changes: ## Should be run from within go main dir, not repo root.
	while true;do \
		{ \
			find internal/ -type f | grep -Ev '\\.(pb|twirp).go'; \
			if [ -d "public" ];then \
				find $(DIR) -mindepth 1 -path $(DIR)/public -prune -o -type f -print; \
			else \
				find $(DIR) -type f; \
			fi; \
		} | entr -r -n -d -s -- "sleep 3;$(MAKE) -C $(DIR) $(TARGETS)"; \
	done
