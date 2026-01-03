SHELL := /bin/bash
CONTAINER := go_inventory_dev

file ?= ./...

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

# Run unit tests inside the dev container (fast)
test:
	docker exec -it $(CONTAINER) /usr/local/go/bin/go test $(file) -v

# Run integration tests inside the dev container (requires test DB)
test-integration:
	docker exec -e TEST_DB_HOST=db -e TEST_DB_PORT=3306 -e TEST_DB_USER=root -e TEST_DB_PASSWORD=root \
		$(CONTAINER) /usr/local/go/bin/go test -tags integration -p 1 $(file) -v

# Run both: unit tests then integration tests (sequential)
test-all: test test-integration
