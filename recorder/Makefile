export POSTGRESQL_DSN ?= postgres://postgres@localhost/website_monitor_test?sslmode=disable

help: # show info about targets
	@grep '^[^#[:space:]].*:' $(MAKEFILE_LIST)

test:
	@docker exec gpagdispo_postgres_1 createdb -U postgres -h localhost website_monitor_test 2> /dev/null || exit 0

	go test -v -race ./...

	@docker exec gpagdispo_postgres_1 dropdb -U postgres -h localhost website_monitor_test 2> /dev/null || exit 0

lint:
	@golangci-lint run -E golint ./... || true
