.PHONY: clean test help
APPS := checker recorder

help: # show info about targets
	@grep '^[^#[:space:]].*:' $(MAKEFILE_LIST)

db-dev: # Create DB
	docker exec gpagdispo_postgres_1 createdb -U postgres -h localhost website_monitor 2> /dev/null || exit 0

# This aim to run a set of containers with the dependencies + setting up the database to develop
start-dev:
	docker-compose -f docker-compose-dev.yaml up --detach
	$(MAKE) db-dev

stop-dev:
	docker-compose -f docker-compose-dev.yaml down

clean: stop-dev

test:
	$(MAKE) -C checker test
	$(MAKE) -C recorder test

integration-test: start-dev
	./e2e-test/run.sh

lint:
	$(MAKE) -C checker lint
	$(MAKE) -C recorder lint
	shellcheck e2e-test/run.sh || true
