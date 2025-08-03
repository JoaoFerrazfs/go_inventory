SHELL := /bin/bash
CONTAINER := go_inventory_dev

# Builda dentro do container
build:
	docker exec -it go_inventory_dev go build -buildvcs=false -o main .

# Roda o main compilado
run:
	docker exec -it $(CONTAINER) ./main

# Builda e roda
start: build run

# Abre um terminal interativo no container
sh:
	docker exec -it $(CONTAINER) bash

# Sobe os containers (modo dev)
up:
	docker compose up -d --build

# Para os containers
stop:
	docker compose down
