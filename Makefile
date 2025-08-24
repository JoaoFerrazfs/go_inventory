SHELL := /bin/bash
CONTAINER := go_inventory_dev

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
