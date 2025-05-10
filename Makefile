.DEFAULT_GOAL := help

.PHONY: help
help: ## Print this help message.
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo ""
	@grep -E '^[a-zA-Z_0-9-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build-image
build-image: ## Build the Docker image.
	docker build -t mcp-nats:latest .

.PHONY: lint 
lint:  ## Lint the Go code.
	golangci-lint run

.PHONY: run
run: ## Run the MCP server in stdio mode.
	go run ./cmd/mcp-nats

.PHONY: run-sse
run-sse: ## Run the MCP server in SSE mode.
	go run ./cmd/mcp-nats --transport sse --log-level debug --debug

.PHONY: run-test-services
run-test-services: ## Run the docker-compose services required for the unit and integration tests.
	docker-compose up -d --build
