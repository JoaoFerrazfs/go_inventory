SHELL := /bin/bash
CONTAINER := go_inventory_dev

up:
	docker compose up -d --build

run:
	docker exec -it $(CONTAINER) ./main

sh:
	docker exec -it $(CONTAINER) bash

stop:
	docker compose down
