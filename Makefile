-include deploy.env

## Set binary commands
GO := go
NODEMON := nodemon

## Set the flags
GOFLAGS := #-mod=vendor
LDFLAGS := -ldflags="-s -w"
GCFLAGS := -gcflags="-m -l"
TESTFLAGS := -v

## Define the source files
ENTRYPOINT := ./tools/test/entrypoint.go

test:
	@$(GO) test $(TESTFLAGS) $(shell $(GO) list ./... | grep -v '/tools')

watch-test-case:
	@$(NODEMON) --watch './**/*.go' --signal SIGTERM --exec $(GO) run $(GOFLAGS) $(LDFLAGS) $(ENTRYPOINT)

run-test-case:
	@$(GO) run $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) $(ENTRYPOINT)
