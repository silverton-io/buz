.PHONY: help build-docker buildx-deploy
S=silverton
REGISTRY:=us-east1-docker.pkg.dev/silverton-io/docker
VERSION:=$(shell cat .VERSION)
TEST_PROFILE=testprofile.out

build-docker: ## Build honeypot docker image locally
	docker build -f deploy/Dockerfile -t honeypot:$(VERSION) .

buildx-deploy: ## Build multi-platform honeypot docker image and push it to internal repo
	docker buildx create --name $(S) || true;
	docker buildx use $(S)
	docker buildx build --platform linux/arm64,linux/amd64 -f deploy/Dockerfile -t $(REGISTRY)/honeypot:$(VERSION)-edge . --push

test-cover-pkg: ## Run tests against pkg, output test profile, and open profile up
	go test ./pkg/... -v -coverprofile=$(TEST_PROFILE) || true
	go tool cover -html=$(TEST_PROFILE) || true

help: ## Display makefile help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
