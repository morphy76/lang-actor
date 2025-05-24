GO := go

GOFLAGS := #-mod=vendor
LDFLAGS := -ldflags="-s -w"
GCFLAGS := -gcflags="-m -l"
TESTFLAGS := -v -count=1 -timeout=2s

test:
	@$(GO) test $(TESTFLAGS) $(shell $(GO) list ./... | grep -vE '/tools/|/examples/')

run-echo-case:
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/actors/echo/run.go

run-pingpong-case:
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/actors/pingpong/run.go

run-selfpingpong-case:
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/actors/selfpingpong/run.go

run-sort-case:
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/actors/sort/run.go

run-calculator-case:
	@echo "With Full Vibes (Github Copilot using Claude 3.7 Sonnet)"
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/actors/calculator/run.go

run-echowithchild-case:
	@echo "With Full Vibes (Github Copilot using Claude 3.7 Sonnet)"
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/actors/echowithchild/run.go

run-simplegraph-case:
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./examples/graph/simple/run.go
