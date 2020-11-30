TAG         ?=  latest

.SUFFIXES:
.PHONY: help build push	test

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build docker image
	@echo Build Balerter $(TAG)
	docker build --build-arg version=$(TAG) -t negasus/ip2location:$(TAG) -f Dockerfile .

push: ## Push docker image
	@echo Push Balerter $(TAG)
	docker push negasus/ip2location:$(TAG)

test: ## Run tests
	go test -mod=vendor ./...
