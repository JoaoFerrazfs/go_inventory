# Testes E2E — Guia Geral do Projeto

Este documento detalha o padrão, arquitetura e dicas para todos os testes E2E do projeto.

---

## Visão Geral

Os testes E2E (end-to-end) automatizados validam o comportamento da aplicação de ponta a ponta, integrando backend, banco de dados, MinIO e outros serviços, em ambiente isolado via Docker Compose. Eles cobrem fluxos reais de autenticação, autorização, regras de negócio e integração entre componentes.

## Padrão e arquitetura dos testes E2E
- **Setup DRY e helpers reutilizáveis:** O setup do banco e seed de dados (ex: usuário admin) é feito por helpers Go (`SupplyInventory/tests/e2e/testsetup`), sem dependência de comandos docker ou scripts externos. Isso garante testes limpos, reprodutíveis e fáceis de expandir.
- **Ambiente isolado:** Cada teste e2e sobe o backend, faz o seed necessário, executa requests reais e valida o comportamento esperado, tudo em containers dedicados.
- **Expansível:** Novos cenários e domínios podem ser testados reutilizando o padrão de setup e helpers.

## Como rodar os testes E2E

**Pré-requisitos:**
- Docker e Docker Compose instalados.
- Arquivo `.env.test` presente na raiz do projeto (já versionado).

**Comando recomendado:**
```bash
docker-compose run --rm test-e2e
```

Esse comando:
- Sobe os containers de banco e MinIO de teste automaticamente.
- Executa todos os testes e2e em um container limpo, com todas as variáveis de ambiente corretas.
- Garante que o backend será iniciado, o seed será feito via código Go, e todos os fluxos serão validados.

## O que os testes E2E cobrem
- Seed automático de dados necessários (ex: admin, usuários, produtos) direto no banco de teste.
- Login, obtenção de token JWT, requests autenticados.
- Validação de acesso permitido e negado para cada role (ex: RBAC), autenticação, regras de negócio, etc.
- Teste de rotas protegidas, integrações e respostas esperadas (200/403/401/4xx/5xx).

## Troubleshooting
- Se houver falha de conexão com o banco, verifique se os containers `db-test` e `minio-test` estão rodando (`docker-compose ps`).
- Confira se o `.env.test` está correto e atualizado.
- Veja logs do backend no output do teste para detalhes de erro.

## Vantagens do padrão adotado
- Setup DRY e reutilizável para todos os testes e2e.
- Sem dependência de comandos docker dentro do teste Go.
- Fácil expansão para novos cenários, domínios e roles.

---

## Expansão
- Para criar novos testes e2e, basta reutilizar o helper de setup/seed.
- Adicione novos cenários de autenticação, autorização, regras de negócio e integrações conforme necessário.
- Consulte exemplos existentes em `SupplyInventory/tests/e2e`.

---

## Exemplos
- **RBAC:** Teste de controle de acesso por função, seed automático do admin, validação de rotas protegidas.
- **Usuários:** Teste de criação, autenticação e fluxo de login.
- **Produtos:** (Exemplo futuro) Teste de cadastro, listagem e regras de estoque.
