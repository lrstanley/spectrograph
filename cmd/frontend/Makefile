#!/usr/bin/make -f

.DEFAULT_GOAL := build

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

fetch: ## Fetches the necessary dependencies to build.
	test -d public/node_modules || npm install

clean: ## Cleans up generated files/folders from the build.
	/bin/rm -rfv "tmp" "dist/*"

clean-cache: ## Cleans up generated cache (speeds up during dev time).
	/bin/rm -rfv ".cache"

build: fetch clean clean-cache ## Builds the frontend.
	npm run build

debug: ## Runs the frontend in debug mode
	npm run watch
