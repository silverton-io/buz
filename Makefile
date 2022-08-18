.PHONY: help run build-docker buildx-deploy test-cover-pkg
S=silverton
REGISTRY:=us-east1-docker.pkg.dev/silverton-io/docker
VERSION:=$(shell cat .VERSION)
HONEYPOT_DIR=./cmd/honeypot/*.go
TEST_PROFILE=testprofile.out

build:
	go build -ldflags="-X main.VERSION=$(VERSION)" -o honeypot $(HONEYPOT_DIR)

run: ## Run honeypot locally
	go run -ldflags="-X 'main.VERSION=x.x.dev'" $(HONEYPOT_DIR)

bootstrap-dev: ## Bootstrap development environment
	curl https://raw.githubusercontent.com/silverton-io/honeypot/main/examples/devel/honeypot/simple.conf.yml -o config.yml;
	make run

bootstrap-dev-destinations: ## Bootstrap various containerized database/stream systems
	docker-compose -f examples/devel/docker-compose.yml up -d

build-docker: ## Build local honeypot image
	docker build -f deploy/Dockerfile -t honeypot:$(VERSION) .

buildx-deploy: ## Build multi-platform honeypot image and push it to edge repo
	docker buildx create --name $(S) || true;
	docker buildx use $(S)
	docker buildx build --platform linux/arm64,linux/amd64 -f deploy/Dockerfile -t $(REGISTRY)/honeypot:$(VERSION)-edge . --push

test-cover-pkg: ## Run tests against pkg, output test profile, and open profile in browser
	go test ./pkg/... -v -coverprofile=$(TEST_PROFILE) || true
	go tool cover -html=$(TEST_PROFILE) || true

help: ## Display makefile help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
