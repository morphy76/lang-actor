GO := go

GOFLAGS := #-mod=vendor
LDFLAGS := -ldflags="-s -w"
GCFLAGS := -gcflags="-m -l"
TESTFLAGS := -v

test:
	@$(GO) test $(TESTFLAGS) $(shell $(GO) list ./... | grep -v '/tools')

run-echo-case:
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./tools/test/echo/run.go

run-pingpong-case:
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./tools/test/pingpong/run.go

run-selfpingpong-case:
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./tools/test/selfpingpong/run.go

run-sort-case:
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./tools/test/sort/run.go
