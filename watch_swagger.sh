#!/bin/bash

# Monitora todos os arquivos na pasta Controllers e roda swag init dentro do container
find ./SupplyInventory/Application/Controllers -type f | entr -r bash -c 'docker compose exec -T app swag init -g main.go'
