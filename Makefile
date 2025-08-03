# Makefile for lang-actor

.DEFAULT_GOAL := help

# Build variables
GO := go
GOFLAGS := #-mod=vendor
LDFLAGS := -ldflags="-s -w"
GCFLAGS := -gcflags="-m -l"
TESTFLAGS := -v -count=1 -timeout=2s

# Declare phony targets
.PHONY: help
help: ## Show this help message
	@echo "Lang-Actor Makefile"
	@echo "==================="
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\033[36m\033[0m"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Testing
.PHONY: test
test: ## Run all tests (excluding tools and examples)
	@$(GO) test $(TESTFLAGS) $(shell $(GO) list ./... | grep -vE '/tools/|/examples/')

##@ Actor Examples
.PHONY: run-echo-case run-pingpong-case run-selfpingpong-case run-sort-case run-calculator-case run-echowithchild-case
run-echo-case: ## Run the echo actor example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/actors/echo/run.go

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
.PHONY: run-simplegraph-case run-forknode-case run-forkgraph-case
run-simplegraph-case: ## Run the simple graph example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/graph/simple/run.go

run-forknode-case: ## Run the forknode graph example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/graph/forknode/run.go

run-forkgraph-case: ## Run the forkgraph graph example
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/graph/forkgraph/run.go
