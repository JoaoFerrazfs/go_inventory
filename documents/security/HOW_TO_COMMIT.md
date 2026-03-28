# 🔧 RBAC - Como Fazer Commit (Git Flow)

Guia para fazer commit da implementação de RBAC seguindo as regras do projeto.

---

## 📋 Antes de Commitar

### Checklist Pre-Commit

```bash
# 1. Verificar status
cd /root/go_inventory
git status

# 2. Rodar testes RBAC
docker exec go_inventory_dev go test ./SupplyInventory/Application/Middlewares -v -run TestRBAC

# 3. Rodar testes de Auth
docker exec go_inventory_dev go test ./SupplyInventory/Application/Controllers/Auth -v

# 4. Rodar todos os testes (opcional, mas recomendado)
docker exec go_inventory_dev go test ./... -v

# 5. Build check
docker exec go_inventory_dev go build ./...
```

Resultado esperado:
```
✅ PASS: TestRBACMiddleware_RequireAdmin_Success
✅ PASS: TestRBACMiddleware_RequireAdmin_Forbidden
✅ PASS: TestRBACMiddleware_RequireRole_MultipleRoles
✅ PASS: TestRBACMiddleware_NoRoleSet
✅ PASS: TestRBACMiddleware_InvalidRoleType
✅ Build OK
```

---

## 📝 Mensagem de Commit

Seguindo o padrão do projeto (ver `.github/instructions/copilot-commit-message.instructions.md`):

### Formato

```
feat: implement RBAC (Role-Based Access Control) middleware

Features:
- Add RBACMiddleware with support for (admin, user) roles
- Add RequireRole(), RequireAdmin(), RequireAny() middleware functions
- Support both OOP and functional programming styles
- Full error handling with type validation

Implementation:
- Create SupplyInventory/Application/Middlewares/rbac.go
- Add RBACMiddleware with granular role checking
- Provide package-level convenience functions
- Return 403 Forbidden for insufficient permissions

Testing:
- Add comprehensive rbac_test.go with 5 test cases
- Test scenarios: success, forbidden, multiple roles, missing role, invalid type
- All tests passing (5/5 ✓)
- Update auth_test.go stubs to match new JWTService interface

Documentation:
- Add RBAC_IMPLEMENTATION_GUIDE.md (11 sections, architecture & implementation)
- Add RBAC_EXAMPLES.md (7 practical examples with curl tests)
- Add RBAC_SUMMARY.md (changes summary & roadmap)
- Add RBAC_NEXT_STEPS.md (checklist for next phases)
- Add RBAC_README_PT-BR.md (Portuguese executive summary)

Migration:
- Backend ready for UserEntity.Role field addition
- Ready to update AuthController.Login() with role assignment
- Ready to add role-based access to new routes

Architecture:
- Role extracted from JWT claims (not client-provided)
- Middleware validates role before handler execution
- Support for multiple allowed roles per route
- Easy to extend with new roles (manager, operator, viewer)

Related: 
- Issue: RBAC Implementation #1
- Ref: PROJECT_REVIEW_2026.md - section 2.1.3

Breaking Changes: None
Migration: None required (backend only)
Database: No changes
```

### Exemplo Simplificado (Se quiser mais curto)

```
feat: implement RBAC middleware with admin and user roles

- Add RBACMiddleware supporting multiple roles
- Implement RequireRole(), RequireAdmin(), RequireAny()
- Add 5 comprehensive tests (all passing)
- Update AuthMiddleware to extract role from JWT
- Update auth test stubs for new JWTService interface
- Add 50+ pages of documentation with examples

Closes: #RBAC-001
```

---

## 🔄 Passo a Passo para Commitar

### 1️⃣ Preparar Mudanças

```bash
cd /root/go_inventory

# Verificar arquivos modificados/criados
git status

# Esperado:
# ✅ Untracked (novos):
#    - SupplyInventory/Application/Middlewares/rbac.go
#    - SupplyInventory/Application/Middlewares/rbac_test.go
#    - documents/security/RBAC_*.md (5 arquivos)
#
# ✅ Modified (modificado):
#    - SupplyInventory/Application/Middlewares/auth_test.go
```

### 2️⃣ Stage Everything

```bash
git add -A

# Verificar
git status

# Deve estar "tudo para commit"
```

### 3️⃣ Rodar Testes Finais

```bash
# Rodar apenas testes RBAC
docker exec go_inventory_dev go test ./SupplyInventory/Application/Middlewares -v -run TestRBAC

# Saída esperada:
# === RUN   TestRBACMiddleware_RequireAdmin_Success
# --- PASS: TestRBACMiddleware_RequireAdmin_Success (0.00s)
# ...
# PASS
# ok      go_inventory/SupplyInventory/Application/Middlewares    0.007s
```

### 4️⃣ Criar Commit

```bash
git commit -m "feat: implement RBAC (Role-Based Access Control) middleware

Features:
- Add RBACMiddleware with support for (admin, user) roles
- Implement RequireRole(), RequireAdmin(), RequireAny() functions
- Support both OOP and functional programming styles
- Full error handling with type validation

Testing:
- Add comprehensive rbac_test.go with 5 test cases
- All tests passing (5/5)
- Update auth_test.go stubs

Documentation:
- Add RBAC_IMPLEMENTATION_GUIDE.md
- Add RBAC_EXAMPLES.md
- Add RBAC_SUMMARY.md
- Add RBAC_NEXT_STEPS.md
- Add RBAC_README_PT-BR.md

Migration:
- Backend ready for UserEntity.Role field
- Ready to update AuthController.Login()
- Ready to add role-based routes"
```

### 5️⃣ Verificar Commit

```bash
# Ver o commit foi criado
git log --oneline -1

# Esperado:
# abc1234 feat: implement RBAC (Role-Based Access Control) middleware

# Ver estrutura completa
git log --pretty=fuller -1

# Ver diferenças
git show --stat

# Esperado: mostra arquivos criados/modificados
```

---

## 🔀 Branches (Recomendado)

Se preferir usar branch (recomendado para PR):

### Criar Branch Feature

```bash
# Criar branch seguindo padrão
git checkout -b feat/rbac-implementation

# Ou
git checkout -b features/rbac

# Fazer as mudanças e commit
git add -A
git commit -m "feat: implement RBAC middleware..."

# Push
git push origin feat/rbac-implementation
```

### No GitHub

Depois criar **Pull Request** de `feat/rbac-implementation` para `main`

**PR Checklist:**
- [ ] Testes passando (CI/CD)
- [ ] Documentação atualizada
- [ ] Sem conflitos com `main`
- [ ] Code review aprovado

---

## 📤 Push para Remoto

### Se Trabalhando com Main

```bash
# Antes de fazer push, trazer atualizações
git pull origin main

# Push
git push origin main

# Verificar no GitHub
# Deve aparecer seu commit no histórico
```

### Se Trabalhando com Branch

```bash
# Push branch para remoto
git push origin feat/rbac-implementation

# No GitHub, abre PR automaticamente
# Ou criar PR manualmente
```

---

## ✅ Verificação Pós-Commit

### Checklist Final

```bash
# 1. Commit foi criado?
git log --oneline -5 | head -1
# Deve mostrar seu commit

# 2. Arquivos corretos no commit?
git show --name-status
# Deve listar os 7 arquivos

# 3. Diff correto?
git show | head -50
# Deve mostrar conteúdo do rbac.go

# 4. Remote updated?
git log origin/main --oneline -1
# Deve mostrar seu commit
```

### No GitHub

1. Acesse: https://github.com/JoaoFerrazfs/go_inventory
2. Vá para **Commits**
3. Procure pelo seu commit `feat: implement RBAC...`
4. Clique para ver detalhes

---

## 🔍 Se Algo Der Errado

### Desfazer Último Commit (sem perder mudanças)

```bash
git reset --soft HEAD~1

# Mudanças voltarão para staging
# Você pode refazer o commit
```

### Desfazer Último Commit (perdendo mudanças)

```bash
# ⚠️ CUIDADO! Irreversível
git reset --hard HEAD~1
```

### Se Fez Push e Quer Desfazer

```bash
# Revert (criar novo commit que desfaz)
git revert HEAD

# Push
git push origin main
```

---

## 📊 Exemplo Completo do Processo

```bash
# 1. Status
$ git status
On branch main
Untracked files:
    SupplyInventory/Application/Middlewares/rbac.go
    SupplyInventory/Application/Middlewares/rbac_test.go
    documents/security/RBAC_*.md

Modified:
    SupplyInventory/Application/Middlewares/auth_test.go

# 2. Stage
$ git add -A
$ git status
Changes to be committed:
    new file: SupplyInventory/Application/Middlewares/rbac.go
    ...

# 3. Testar
$ docker exec go_inventory_dev go test ./SupplyInventory/Application/Middlewares -v -run TestRBAC
PASS
ok  go_inventory/SupplyInventory/Application/Middlewares    0.007s

# 4. Commit
$ git commit -m "feat: implement RBAC..."
[main 234abe5] feat: implement RBAC...
 7 files changed, 500 insertions(+), 5 deletions(-)

# 5. Verify
$ git log --oneline -1
234abe5 feat: implement RBAC (Role-Based Access Control) middleware

# 6. Push
$ git push origin main
Enumerating objects: 12, done.
Counting objects: 100% (12/12), done.
Writing objects: 100% (10/10), 2.50 KiB | 2.50 MiB/s, done.
To github.com:JoaoFerrazfs/go_inventory.git
   5d4e3c2..234abe5  main -> main
```

---

## 🎯 Comando Completo (Copy-Paste)

Para facilitar, aqui está o comando completo:

```bash
cd /root/go_inventory && \
git add -A && \
docker exec go_inventory_dev go test ./SupplyInventory/Application/Middlewares -v -run TestRBAC && \
git commit -m "feat: implement RBAC (Role-Based Access Control) middleware

Features:
- Add RBACMiddleware with (admin, user) role support
- Implement RequireRole(), RequireAdmin(), RequireAny() middleware
- Full error handling with role type validation

Testing:
- Add rbac_test.go with 5 comprehensive test cases (all passing)
- Update auth_test.go stubs for JWTService interface

Documentation:
- Add 5 markdown files (~50 pages) with guides and examples
- RBAC_IMPLEMENTATION_GUIDE.md - architecture
- RBAC_EXAMPLES.md - 7 practical examples
- RBAC_README_PT-BR.md - Portuguese summary
- RBAC_SUMMARY.md - changes overview
- RBAC_NEXT_STEPS.md - checklist

Backend ready for UserEntity.Role integration." && \
git log --oneline -1
```

---

## 📌 Notas Importantes

### Padrões do Projeto

Verificar `.github/instructions/`:

1. **copilot-commit-message.instructions.md** - Formato de mensagem
2. **copilot-pull-request.instructions.md** - Padrão de PR
3. **tests-patters.instructions.md** - Padrão de testes

### Antes de Commitar

- ✅ Testes passando
- ✅ Sem conflitos
- ✅ Código formatado
- ✅ Mensagem clara
- ✅ Referência a issues (se houver)

### Após Commitar

- ✅ Verificar no GitHub
- ✅ CI/CD passou
- ✅ Code review (se necessário)
- ✅ Documentação atualizada

---

## 🚀 Próximas Etapas Após o Commit

1. **Comunicar a mudança:**
   - Atualizar card no Kanban
   - Notificar time (Discord/Slack)
   - Adicionar link do commit

2. **Próxima tarefa:**
   - Iniciar implementation de UserEntity.Role
   - Ver `documents/security/RBAC_NEXT_STEPS.md`

3. **PR e Code Review:**
   - Se em branch, abrir PR
   - Pedir revisão dos colegas
   - Fazer ajustes se necessário

---

## 📞 Dúvidas no Git?

### Comandos Úteis

```bash
# Ver histórico
git log --oneline -10

# Ver mudanças não commited
git diff

# Ver mudanças commited
git diff --cached

# Ver um commit específico
git show HASH

# Revertendo coisas
git reset --soft HEAD~1        # Desfaz mantendo mudanças
git revert HASH                # Cria commit que desfaz

# Branches
git branch -a                  # Lista todas
git checkout -b feature/name   # Criar nova branch
git push origin feature/name   # Push da branch
```

---

**Pronto para fazer commit? Siga o guia acima e boa sorte! 🚀**

