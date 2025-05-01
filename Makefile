GO := go

GOFLAGS := #-mod=vendor
LDFLAGS := -ldflags="-s -w"
GCFLAGS := -gcflags="-m -l"
TESTFLAGS := -v

test:
	@$(GO) test $(TESTFLAGS) $(shell $(GO) list ./... | grep -v '/tools')

run-echo-case:
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) ./tools/test/entrypoint.go
