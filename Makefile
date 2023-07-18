# ---------------------------------------- COMMON START ----------------------------------------------------------------
.PHONY: test
test:
	CGO_ENABLED=1 go test -v -race -count=1 ./...

.PHONY: lint
lint:
	@golangci-lint run
# ---------------------------------------- COMMON END ------------------------------------------------------------------