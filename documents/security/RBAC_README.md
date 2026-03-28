# RBAC (Controle de Acesso por Função) — Guia Prático

## Decisão de Arquitetura

- **RBAC Middleware**: Permite proteger rotas por função (role) de usuário.
- **Roles padrão**: `admin` (acesso total) e `user` (acesso limitado).
- **Extensível**: Fácil adicionar novas roles (ex: manager, operator).
- **JWT**: O token inclui a role do usuário, propagada para o contexto do request.

---

## Como Proteger Rotas

### Apenas Admin

```go
adminGroup := apiV1.Group("/admin/users")
adminGroup.Use(authMiddleware.Handler())      // Autenticação
adminGroup.Use(middlewares.RequireAdmin())    // Só admin acessa
userController.RegisterAdminRoutes(adminGroup)
```

### Múltiplas Roles Permitidas

```go
managerGroup := apiV1.Group("/manager/reports")
managerGroup.Use(authMiddleware.Handler())
managerGroup.Use(middlewares.RequireRole(
    constants.AdminRole,
    // constants.ManagerRole, // quando implementar
))
reportController.Register(managerGroup)
```

---

## Exemplo Real de Uso via Terminal

### 1. Admin acessando rota protegida

```bash
curl -X GET "http://localhost:8000/api/v1/admin/racks?page=1" \
  -H "Authorization: Bearer <TOKEN_ADMIN>"
```
**Resposta esperada:** 200 OK, dados retornados.

### 2. User tentando acessar rota admin

```bash
curl -X GET "http://localhost:8000/api/v1/admin/racks?page=1" \
  -H "Authorization: Bearer <TOKEN_USER>"
```
**Resposta esperada:** 403 Forbidden
```json
{ "error": "insufficient permissions" }
```

### 3. Rota com múltiplas roles (admin e manager)

```bash
curl -X GET "http://localhost:8000/api/v1/manager/reports/summary" \
  -H "Authorization: Bearer <TOKEN_ADMIN>"
# ou
curl -X GET "http://localhost:8000/api/v1/manager/reports/summary" \
  -H "Authorization: Bearer <TOKEN_MANAGER>"
```
**Resposta esperada:** 200 OK para admin/manager, 403 para user comum.

---

## Como funciona internamente

- O middleware lê a role do JWT e coloca no contexto.
- O RBAC verifica se a role do usuário está entre as permitidas.
- Se não estiver, retorna 403 Forbidden.

---

## Testes

Todos os testes automatizados cobrem:
- Acesso permitido e negado para cada role.
- Cenários de erro (role ausente ou inválida).

Para rodar:
```bash
docker exec go_inventory_dev go test ./SupplyInventory/Application/Middlewares -v
```

---

## Checklist para expandir

- Adicione novas roles em `SupplyInventory/Domain/constants/roles.go`
- Use `middlewares.RequireRole(constants.NovaRole)` para proteger rotas.
- Garanta que o JWT inclua a nova role ao autenticar.
