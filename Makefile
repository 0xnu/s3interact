VERSION=1.1.2
COMMIT_MSG ?= "format options"

execute: ## Execute Locally
	@go mod init s3interact
	@go mod tidy
	@go run .

build: ## Build Single Binary for Local OS
	@go build -v ./

package: ## Build for Multi OS (linux 386, amd64).
	@chmod +x package.sh && ./package.sh

acp: ## Add, Commit and Push
	@git add .
	@git commit -s -m $(COMMIT_MSG)
	@git push

tag: ## Tag
	@git tag "v$(VERSION)"
	@git push --tag

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\\033[36m%-30s\\033[0m %s\\n", $$1, $$2}'

.PHONY: help

.DEFAULT_GOAL := help
