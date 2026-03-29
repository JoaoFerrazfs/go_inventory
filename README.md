# Go Inventory System

Sistema de gerenciamento de paletes, racks e produtos desenvolvido em Go, utilizando Gin, GORM e Docker.

## 🚀 Tecnologias e Bibliotecas

- **Linguagem:** Go 1.24  
- **Framework Web:** Gin  
- **ORM:** GORM com MySQL  
- **Gerenciamento de Dependências:** Go Modules  
- **Injeção de Dependências:** dig (`go.uber.org/dig`)  
- **QRCode:** go-qrcode (`github.com/skip2/go-qrcode`)  )  
- **Docker:** Contêinerização da aplicação e do banco de dados  
- **Banco de dados:** MySQL  

## 📋 Visão Geral do Projeto

Este é um sistema de gerenciamento de inventário baseado em Go para manipulação de paletes, racks de paletes e produtos paletizados. Ele fornece APIs RESTful para operações CRUD, autenticação de usuários, geração de códigos QR para paletes e integração com banco de dados MySQL. O projeto enfatiza arquitetura limpa, injeção de dependências e conteinerização com Docker. 🏗️📦

## 🏛️ Arquitetura

O código segue uma arquitetura em camadas inspirada na Clean Architecture:
- **Camada de Domínio**: Contém entidades definindo modelos de negócio e regras. 🎯
- **Camada de Infraestrutura**: Gerencia persistência de dados com repositórios e conexões de banco. 🗄️
- **Camada de Aplicação**: Inclui serviços, controladores e middlewares para lógica e manipulação de API. ⚙️
- **Injeção de Dependências**: Usa `dig` para conectar componentes, promovendo testabilidade e modularidade. 🔗
- **Ponto de Entrada**: Configura o roteador Gin, CORS, arquivos estáticos, migração de banco e registro de rotas. 🚀

## ✨ Funcionalidades Principais

- **Autenticação e Autorização**: Login/refresh baseado em JWT. Usuários gerenciados com serviços dedicados. 🔐
- **Operações CRUD**:
  - Paletes: Criar, listar, encontrar por ID, atualizar, deletar (com geração de QR code). 📦
  - Racks de Paletes: CRUD similar, com rastreamento de capacidade e porcentagem de uso. 🗂️
  - Produtos Paletizados: Adicionar/remover produtos de paletes. 📋
- **Integração com QR Code**: Gera códigos QR linkando para detalhes de paletes, armazenados e servidos estaticamente. 📱
- **Documentação de API**: Integração com Swagger com auto-geração. 📖
- **Tratamento de Erros**: Erros personalizados para respostas consistentes. ⚠️
- **Validação**: Estruturas de requisição com binding e erros formatados. ✅

## 🛠️ Tecnologias e Ferramentas

- **Linguagem**: Go 1.24 🐹
- **Framework Web**: Gin para roteamento e middleware. 🌐
- **ORM**: GORM com driver MySQL. 🗃️
- **Banco de Dados**: MySQL, com migrações. 💾
- **Conteinerização**: Docker Compose para app e DB. 🐳
- **Outras Bibliotecas**: `go-qrcode` para geração de QR, `dig` para DI, `gin-swagger` para docs, `godotenv` para vars de ambiente. 📚
- **Ferramentas de Build/Dev**: `Makefile` para comandos como `make dev` para executar containers e monitorar Swagger. 🔧

## 🔬 Testing

We follow the project's testing standards (see `documents/Testing/TestingStandards.md`). Quick commands:

Run unit tests inside the dev container:

```bash
docker exec -it go_inventory_dev /usr/local/go/bin/go test ./... -v
```

Run integration tests inside the dev container (example using MySQL container IP):

```bash
docker exec -e TEST_DB_HOST=172.17.0.2 -e TEST_DB_PORT=3306 -e TEST_DB_USER=root -e TEST_DB_PASSWORD=root go_inventory_dev /usr/local/go/bin/go test -tags integration ./... -v
```

Branch and commit policy for tests: create a focused branch named `test/<scope>` and commit test files only after verifying they pass locally. Use `test:` in the commit message prefix for test-only commits.

## 🧪 Testes: Unitários, Integração e E2E

O projeto utiliza build tags para separar os tipos de teste:

- **Unitários**: `//go:build unit` — Não dependem de banco ou serviços externos.
- **Integração**: `//go:build integration` — Usam banco real, testam integração entre componentes.
- **E2E**: `//go:build e2e` — Simulam uso real da API, sobem o backend e testam via HTTP.

### Como rodar os testes

Todos os comandos devem ser executados na raiz do projeto:

- **Testes unitários:**
  ```sh
  make test-unit
  ```
- **Testes de integração:**
  ```sh
  make test-integration
  ```
- **Testes E2E:**
  ```sh
  make test-e2e
  ```
- **Todos os testes (sequencial):**
  ```sh
  make test-all
  ```

> Os testes são executados dentro do container Docker `go_inventory_dev` e usam o pretty reporter customizado.
> Para integração/E2E, o banco de testes deve estar disponível (veja docker-compose).

### Estrutura dos arquivos de teste
- Testes unitários: `*_test.go` com tag `unit`.
- Testes de integração: `*_integration_test.go` com tag `integration`.
- Testes E2E: `SupplyInventory/tests/e2e/rbac_e2e_test.go` com tag `e2e`.

Consulte os comentários nos próprios arquivos para exemplos e padrões de escrita.
