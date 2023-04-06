.PHONY: help run bootstrap-destinations build-docker buildx-deploy test-cover-pkg
S=silverton
REGISTRY:=us-east1-docker.pkg.dev/silverton-io/docker
VERSION:=$(shell cat .VERSION)
BUZ_DIR="./cmd/buz/"
TEST_PROFILE=testprofile.out

build:
	go build -ldflags="-X main.VERSION=$(VERSION)" -o build/buz $(BUZ_DIR)

run: ## Run buz locally
	go run -ldflags="-X 'main.VERSION=x.x.dev'" $(BUZ_DIR)

debug: ## Run buz locally with debug
	DEBUG=1 go run -ldflags="-X 'main.VERSION=x.x.dev'" $(BUZ_DIR)

bootstrap: ## Bootstrap development environment
	test -f config.yml || cp ./examples/devel/buz/simple.conf.yml config.yml;
	make debug

bootstrap-destinations: ## Bootstrap various containerized database/stream systems
	docker-compose -f examples/devel/docker-compose.yml up -d

build-docker: ## Build local buz image
	docker build -f deploy/Dockerfile -t buz:$(VERSION) .

buildx-deploy: ## Build multi-platform buz image and push it to edge repo
	docker buildx create --name $(S) || true;
	docker buildx use $(S)
	docker buildx build --platform linux/arm64,linux/amd64 -f deploy/Dockerfile -t $(REGISTRY)/buz:$(VERSION)-edge . --push

lint: ## Lint go code
	@golangci-lint run --config .golangci.yml

test: ## Run tests against pkg
	@go test ./pkg/...

test-cover-pkg: ## Run tests against pkg, output test profile, and open profile in browser
	go test ./pkg/... -v -coverprofile=$(TEST_PROFILE) || true
	go tool cover -html=$(TEST_PROFILE) || true

help: ## Display makefile help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
