#!/usr/bin/make -f

.DEFAULT_GOAL := build

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

fetch: ## Fetches the necessary dependencies to build.
	npm install --no-audit --no-fund --loglevel verbose

clean: ## Cleans up generated files/folders from the build.
	/bin/rm -rfv "tmp" "dist/*"

clean-cache: ## Cleans up generated cache (speeds up during dev time).
	/bin/rm -rfv ".cache"

build: fetch clean clean-cache ## Builds the frontend.
	./node_modules/parcel-bundler/bin/cli.js build \
		--no-autoinstall \
		--no-source-maps \
		--log-level 4 \
		--target browser \
		--public-url '/static/dist' \
		src/index.html

debug: fetch clean ## Runs the frontend in debug mode
	./node_modules/parcel-bundler/bin/cli.js watch \
		--hmr-port 8081 \
		--log-level 4 \
		--target browser \
		--public-url '/static/dist' \
		src/index.html
