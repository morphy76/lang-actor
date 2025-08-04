# Makefile for lang-actor

.DEFAULT_GOAL := help

# Build variables
GO := go
GOFLAGS := #-mod=vendor
LDFLAGS := -ldflags="-s -w"
GCFLAGS := -gcflags="-m -l"
TESTFLAGS := -v -count=1 -timeout=30s -race -failfast -shuffle=on -coverprofile=coverage.out
LINTFLAGS := #-v
PACKAGES := $(shell $(GO) list ./... | grep -vE '/tools/|/examples/')

# Declare phony targets
.PHONY: help
help: ## Show this help message
	@echo "Lang-Actor Makefile"
	@echo "==================="
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\033[36m\033[0m"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Static Analysis
.PHONY: lint
lint: ## Run static analysis tools (golint, go vet, etc.)
	@echo "Running static analysis..."
	@$(GO) vet $(LINTFLAGS) $(PACKAGES)

##@ Testing
.PHONY: test
test: lint ## Run all tests with race detection and comprehensive flags
	@$(GO) test $(TESTFLAGS) $(PACKAGES)

.PHONY: test-bench
test-bench: ## Run benchmark tests
	@$(GO) test -v -bench=. -benchmem -timeout=60s $(PACKAGES)

##@ Cleanup
.PHONY: clean clean-test
clean: clean-test ## Clean all generated files

clean-test: ## Clean test artifacts (coverage files, etc.)
	@echo "Cleaning test artifacts..."
	@rm -f coverage.out coverage.html
	@echo "✅ Test artifacts cleaned"

##@ Actor Examples
.PHONY: run-echo-case run-pingpong-case run-selfpingpong-case run-sort-case run-calculator-case run-echowithchild-case run-counter-case
run-echo-case: ## Run the echo actor example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/actors/echo/run.go

run-counter-case: ## Run the counter actor example (from README)
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/actors/counter/run.go

run-pingpong-case: ## Run the pingpong actor example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/actors/pingpong/run.go

run-selfpingpong-case: ## Run the selfpingpong actor example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/actors/selfpingpong/run.go

run-sort-case: ## Run the sort actor example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/actors/sort/run.go

run-calculator-case: ## Run the calculator actor example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/actors/calculator/run.go

run-echowithchild-case: ## Run the echowithchild actor example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/actors/echowithchild/run.go

##@ Graph Examples
.PHONY: run-simplegraph-case run-forknode-case run-forkgraph-case run-loremstream-case run-simple-ollama-case
run-simplegraph-case: ## Run the simple graph example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/graph/simple/run.go

run-forknode-case: ## Run the forknode graph example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/graph/forknode/run.go

run-forkgraph-case: ## Run the forkgraph graph example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/graph/forkgraph/run.go

run-loremstream-case: ## Run the lorem ipsum streaming graph example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/graph/loremstream/run.go

run-simple-ollama-case: ## Run the simple ollama graph example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/graph/simple-ollama/run.go
