SHELL := /bin/bash
CONTAINER := go_inventory_dev

file ?= ./...
REPORTER_ARGS ?=

ifeq ($(strip $(LIST)),)
	# LIST not set
else
	REPORTER_ARGS := -h
	SHOW_REPORTER_HELP := $(LIST)
endif

.PHONY: LIST
LIST: .build-go-reporter
	docker exec $(CONTAINER) sh -c '/tmp/pretty_print_tests -h'

# Start containers normalmente
up:
	docker compose up

# Gera swagger automaticamente ao salvar arquivos
watch-swagger:
	./watch_swagger.sh

# Roda docker + watch do swagger ao mesmo tempo
dev:
	$(MAKE) up &
	$(MAKE) watch-swagger &
	wait


# Build the Go reporter inside the container
.build-go-reporter:
	@echo "Building go test pretty reporter inside container..."
	@docker exec $(CONTAINER) /bin/sh -c 'mkdir -p /tmp/reporter'
	@docker cp scripts/pretty_print_tests.go $(CONTAINER):/tmp/reporter/pretty_print_tests.go
	@docker exec $(CONTAINER) /bin/sh -c 'cd /tmp/reporter && /usr/local/go/bin/go build -o /tmp/pretty_print_tests pretty_print_tests.go'

# Run unit tests inside the dev container (fast) and pretty-print failures inside container
test: .build-go-reporter
	@sh -c 'if [ -n """$(SHOW_REPORTER_HELP)""" ]; then \
		docker exec $(CONTAINER) sh -c "/tmp/pretty_print_tests -h"; \
	else \
		docker exec $(CONTAINER) sh -c "/usr/local/go/bin/go test $(file) -json | /tmp/pretty_print_tests $(REPORTER_ARGS)"; \
	fi'

# Run integration tests inside the dev container (requires test DB)
test-integration: .build-go-reporter
	docker exec -e TEST_DB_HOST=db -e TEST_DB_PORT=3306 -e TEST_DB_USER=root -e TEST_DB_PASSWORD=root $(CONTAINER) \
    sh -c '/usr/local/go/bin/go test -tags integration -p 1 $(file) -json | /tmp/pretty_print_tests $(REPORTER_ARGS)'

# Run both: unit tests then integration tests (sequential)
test-all: test test-integration

# Run tests and pretty-print failures using a JSON parser
test-pretty: .build-go-reporter
	docker exec $(CONTAINER) sh -c '/usr/local/go/bin/go test $(file) -json | /tmp/pretty_print_tests $(REPORTER_ARGS)'

# Testes Go: targets rápidos
#
# make test-unit        # Roda só testes unitários (tag unit)
# make test-integration # Roda só testes de integração (tag integration)
# make test-e2e         # Roda só testes E2E (tag e2e)
# make test-all         # Roda todos (unit, integration, e2e em sequência)

# Testes unitários (rápidos, sem dependências externas)
test-unit: .build-go-reporter
	docker exec $(CONTAINER) sh -c '/usr/local/go/bin/go test -tags unit -p 1 ./... -json | /tmp/pretty_print_tests $(REPORTER_ARGS)'

# Testes de integração (usam banco real, dependências externas)
test-integration: .build-go-reporter
	docker exec -e TEST_DB_HOST=db -e TEST_DB_PORT=3306 -e TEST_DB_USER=root -e TEST_DB_PASSWORD=root $(CONTAINER) \
    sh -c '/usr/local/go/bin/go test -tags integration -p 1 ./... -json | /tmp/pretty_print_tests $(REPORTER_ARGS)'


# Testes end-to-end (E2E, simulam uso real via HTTP)
test-e2e:
	docker-compose run --rm test-e2e

# Roda todos os testes (unit, integration, e2e)
test-all: test-unit test-integration test-e2e
