help: # show info about targets
	@grep '^[^#[:space:]].*:' $(MAKEFILE_LIST)

test:
	@go test -v -race ./...

lint:
	@golangci-lint run -E golint ./... || true
