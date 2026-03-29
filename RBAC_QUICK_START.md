# 🚀 RBAC Implementation - Quick Start

## O que foi implementado?

✅ Sistema completo de **RBAC (Role-Based Access Control)** com roles **admin** e **user**

## Como começar em 5 minutos

### 1️⃣ Dev Environment (Novo Banco)

Se você está iniciando com banco novo, tudo funciona automaticamente:

```bash
# Iniciar Docker
docker-compose up

# O banco é criado e migração roda automaticamente
# UserEntity já terá coluna 'role' com default 'user'
```

✅ **Pronto!** Novo usuários terão `role: "user"` por padrão.

---

### 2️⃣ Banco Existente - Adicionar Coluna Role

Se você tem banco existente, execute **uma vez**:

```sql
-- Conexão ao MySQL
mysql -h db -u root -proot inventory

-- Adicionar coluna role
ALTER TABLE user_entities 
ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user' AFTER password;

-- Criar index (opcional mas recomendado)
CREATE INDEX idx_user_role ON user_entities(role);

-- Verificar
SELECT id, email, role FROM user_entities LIMIT 5;
```

✅ **Pronto!** Todos usuários existentes agora são `role: "user"`.

---

### 3️⃣ Criar Usuário Admin

**Via SQL (RECOMENDADO POR ENQUANTO):**

```sql
-- Promover um usuário existente a admin
UPDATE user_entities 
SET role = 'admin' 
WHERE email = 'admin@example.com';

-- Ou criar novo (com password hasheado)
-- ATENÇÃO: Substituir [HASH_AQUI] por bcrypt hash real
INSERT INTO user_entities (name, email, password, role, created_at, updated_at)
VALUES ('Admin User', 'admin@test.com', '$2a$10$...bcrypt_hash_aqui...', 'admin', NOW(), NOW());
```

**Para gerar hash bcrypt rapidamente:**

```bash
# Opção 1: Usar online (NÃO em produção!)
# https://bcrypt.online/

# Opção 2: Usar Go no container
docker exec go_inventory_dev go run -
# Cole este código:
package main
import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)
func main() {
    hash, _ := bcrypt.GenerateFromPassword([]byte("sua_senha_aqui"), 10)
    fmt.Println(string(hash))
}
```

✅ **Pronto!** Agora tem um admin.

---

### 4️⃣ Testar RBAC

#### A. Login como User Regular

```bash
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "senha123456"
  }'
```

Resposta:
```json
{
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGc...",
  "refresh_token": "eyJ0eXAiOiJKV1QiLCJhbGc..."
}
```

#### B. Tentar Acessar Admin Endpoint (vai falhar)

```bash
curl -X GET http://localhost:8000/api/v1/admin/racks \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGc..."
```

Resposta:
```json
{
  "error": "insufficient permissions"
}
```

**Status:** `403 Forbidden` ✅

#### C. Login como Admin

```bash
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "senha123456"
  }'
```

#### D. Acessar Admin Endpoint (vai funcionar)

```bash
curl -X GET http://localhost:8000/api/v1/admin/racks \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGc..."
```

Resposta:
```json
[
  {"id": 1, "name": "Rack A", ...},
  {"id": 2, "name": "Rack B", ...}
]
```

**Status:** `200 OK` ✅

---

## Estrutura de Routes Protegidas

### 🟢 Rotas PÚBLICAS (Sem auth)
```
POST   /api/v1/auth/login              # Login (qualquer um)
```

### 🔵 Rotas AUTENTICADAS (Qualquer role)
```
GET    /api/v1/inventories             # User ou Admin
GET    /api/v1/pallets                 # User ou Admin
PATCH  /api/v1/pallet/products/:id     # User ou Admin
```

### 🔴 Rotas ADMIN ONLY
```
GET    /api/v1/admin/racks             # Apenas ADMIN
```

---

## Token JWT (Como Funciona)

### User Token (após login)
```json
{
  "userID": 5,
  "username": "user@example.com",
  "role": "user",
  "tokenType": "access",
  "exp": 1234567890
}
```

### Admin Token (após login)
```json
{
  "userID": 1,
  "username": "admin@example.com",
  "role": "admin",
  "tokenType": "access",
  "exp": 1234567890
}
```

Middleware RBAC lê `"role"` do token e valida antes de executar endpoint.

---

## Próximos Passos (Fase 2)

### Adicionar Endpoint para Criar Admin via API

```
POST /api/v1/admin/users
Headers: Authorization: Bearer [ADMIN_TOKEN]
Body: {
  "name": "New Admin",
  "email": "newadmin@example.com",
  "password": "senha123456",
  "role": "admin"
}
```

### Adicionar Novas Roles (e.g., Manager)

1. Editar `SupplyInventory/Domain/constants/roles.go`
2. Adicionar `ManagerRole = "manager"`
3. Usarem rotas conforme necessário
4. **Sem mudanças no middleware!**

---

## Troubleshooting

| Problema | Solução |
|----------|---------|
| ❌ "user role not found in token" | Fazer novo login (token antigo não tem role) |
| ❌ `403 Forbidden` em admin endpoint | Verificar se user tem role=admin no BD |
| ❌ "Invalid token format" | Verificar coluna `role` foi adicionada no BD |
| ✅ Tudo ok mas não funciona | Ver logs: `docker logs go_inventory_dev` |

---

## Arquivos para Consulta

| Arquivo | O que tem |
|---------|-----------|
| `documents/security/RBAC_GUIDE.md` | Guia completo + exemplos código |
| `documents/security/RBAC_IMPLEMENTATION.md` | Resumo implementação + file changes |
| `SupplyInventory/Domain/constants/roles.go` | Definição de roles |
| `SupplyInventory/Application/Middlewares/rbac.go` | Middleware de RBAC |
| `SupplyInventory/Application/Middlewares/rbac_test.go` | Testes do RBAC |

---

## Checklist

- [x] Roles criadas (admin, user)
- [x] UserEntity com field role
- [x] JWT inclui role nos claims
- [x] AuthMiddleware extrai role
- [x] RBACMiddleware criado
- [x] Rotas admin protegidas
- [x] Documentação completa
- [x] Testes de RBAC
- [ ] **PRÓXIMO:** Implementar endpoint para criar/editar usuário com role (REST API)
- [ ] **PRÓXIMO:** Adicionar novas roles conforme necessidade

---

**Status:** ✅ **PRONTO PARA USAR**

