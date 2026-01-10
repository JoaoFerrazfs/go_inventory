---
applyTo: '**/*_test.go'
---

## Package dos Testes Go

- Sempre utilize um package de teste separado, com sufixo _test (ex: controllers_test, repositories_test), para todos os arquivos de teste Go (_test.go).
- Isso garante que os testes só acessem a API pública do package, promovendo isolamento e simulando o uso real do package por consumidores externos.
- O uso de package separado evita dependências acidentais em funções, variáveis ou tipos não exportados.
- Exemplo:
  - Código: package controllers
  - Teste: package controllers_test

## Testes
- Os testes devem ser fáceis de entender, bem estruturados e organizados.
- Não use siglas; prefira palavras completas em nomes de variáveis, funções, métodos e arquivos de teste.
- Prefira legibilidade simples em vez de soluções excessivamente complexas.
- Para executar testes Go, utilize sempre o container Docker da aplicação.
- Comandos principais (use o binário `go` absoluto dentro do container para maior confiabilidade):
  - Testes unitários (com mocks):
    - docker exec -it go_inventory_dev /usr/local/go/bin/go test ./SupplyInventory/Application/Controllers/Auth/... -v
  - Testes de integração (com banco real):
    - docker exec -e TEST_DB_HOST=db -e TEST_DB_PORT=3306 -e TEST_DB_USER=root -e TEST_DB_PASSWORD=root \
      go_inventory_dev /usr/local/go/bin/go test -tags integration ./SupplyInventory/Application/Controllers/Auth/... -v
  - Todos os testes unitários:
    - docker exec -it go_inventory_dev /usr/local/go/bin/go test ./... -v
  - Todos os testes de integração:
    - docker exec -e TEST_DB_HOST=db -e TEST_DB_PORT=3306 -e TEST_DB_USER=root -e TEST_DB_PASSWORD=root \
      go_inventory_dev /usr/local/go/bin/go test -tags integration ./... -v

  - Executar ambos (unitários e integração) em sequência:
    - docker exec -it go_inventory_dev /usr/local/go/bin/go test ./... -v && \
      docker exec -e TEST_DB_HOST=db -e TEST_DB_PORT=3306 -e TEST_DB_USER=root -e TEST_DB_PASSWORD=root \
      go_inventory_dev /usr/local/go/bin/go test -tags integration ./... -v

- Certifique-se de instalar as dependências de teste no ambiente do container quando necessário:
  - github.com/stretchr/testify
  - github.com/davecgh/go-spew/spew
  - github.com/pmezard/go-difflib/difflib
  - github.com/stretchr/objx
- Use `testify` mocks para simular serviços e repositórios nos testes unitários.
- Para testes de integração, use MySQL real via containers.
- Siga a estrutura de pastas de testes alinhada ao domínio do código.
- Write unit tests for individual components and functions.
- Implement integration tests to verify interactions between components.
  - Each API endpoint must have an integration test.
- Use end-to-end testing para simular cenários reais do usuário.
- Busque cobertura de testes alta, priorizando áreas críticas e complexas.
- Utilize frameworks e ferramentas adequadas ao ecossistema Go.
- Escreva casos de teste claros e descritivos com comportamento esperado.
- Realize testes de carga e performance quando necessário.
- Todos os arquivos de teste Go (_test.go) devem ficar na mesma pasta do arquivo que testam.
  - Exemplo: se existe /SupplyInventory/Application/Controllers/AuthController.go, o teste deve ser /SupplyInventory/Application/Controllers/AuthController_test.go.
  - Isso garante que o coverage funcione corretamente e segue o padrão da comunidade Go.

- **Estrutura dos Testes Unitários:**
  - Use mocks para repositórios e serviços externos.
  - Organize os testes em blocos separados e comentados, seguindo a ordem:
    - // Set
    - // Expectations (se houver)
    - // Actions
    - // Assertions
  - Deixe uma linha em branco entre cada bloco.
  - O bloco `// Expectations` só aparece quando houver expectativas de mocks; caso contrário, omita.
  - Todos os comentários no código de teste devem estar em English.

  **Exemplo ERRADO:**
  ```go
  func TestCreatePallet_Success(t *testing.T) {
      // Set
      palletRepo := new(mocks.PalletRepository)
      palletService := domain.NewPalletService(palletRepo)
      pallet := &domain.Pallet{ID: "123", Name: "Test Pallet"}
      palletRepo.On("Save", pallet).Return(nil)
      // Actions
      err := palletService.CreatePallet(pallet)
      // Assertions
      assert.NoError(t, err)
      palletRepo.AssertExpectations(t)
  }
  ```

  **Exemplo CERTO:**
  ```go
  func TestCreatePallet_Success(t *testing.T) {
      // Set
      palletRepo := new(mocks.PalletRepository)
      palletService := domain.NewPalletService(palletRepo)
      pallet := &domain.Pallet{ID: "123", Name: "Test Pallet"}

      // Expectations
      palletRepo.On("Save", pallet).Return(nil)

      // Actions
      err := palletService.CreatePallet(pallet)

      // Assertions
      assert.NoError(t, err)
      palletRepo.AssertExpectations(t)
  }
  ```

- **Estrutura dos Testes de Integração:**
  - Use banco de dados MySQL real (inventory_test) com AutoMigrate.
  - Utilize transações para isolamento: cada teste roda dentro de uma transação que é revertida ao final.
  - Trunque as tabelas antes de cada teste para garantir estado limpo.
  - Use fixtures para criar dados de teste (CreateTestUser, CreateTestPallet, etc.).
  - Os testes de integração devem usar build tags: //go:build integration
  - Organize os testes em blocos separados e comentados, seguindo a ordem:
    - // Set
    - // Actions
    - // Assertions
  - Use IntegrationTestHelper para setup do banco e routers.
  - Cada teste deve truncar tabelas fora da transação e usar transação para rollback.

  **Exemplo de Teste de Integração:**
  ```go
  //go:build integration
  // +build integration

  package controllers_test

  import (
      "bytes"
      "encoding/json"
      "net/http"
      "net/http/httptest"
      "testing"

      "github.com/stretchr/testify/assert"
      "gorm.io/gorm"

      userRequests "go_inventory/SupplyInventory/Application/Requests/User"
      integration "go_inventory/SupplyInventory/tests/integration"
  )

  func TestIntegration_CreateUser(t *testing.T) {
      h := integration.NewIntegrationTestHelper()
      h.TruncateTables(h.DB)
      h.DB.Transaction(func(tx *gorm.DB) error {
          // Set
          r := h.SetupRouterForUser(tx)

          createReq := userRequests.UserRequest{
              Name:     "Admin",
              Email:    "admin@example.com",
              Password: "admin123",
          }
          body, _ := json.Marshal(createReq)

          // Actions
          w := httptest.NewRecorder()
          req, _ := http.NewRequest("POST", "/api/v1/users/create", bytes.NewBuffer(body))
          req.Header.Set("Content-Type", "application/json")
          r.ServeHTTP(w, req)

          // Assertions
          assert.Equal(t, http.StatusCreated, w.Code)
          return nil
      })
  }
  ```

- **Utilitários de Teste:**
  - SupplyInventory/tests/testutils/db.go: SetupTestDB para criar DB de teste e AutoMigrate.
  - SupplyInventory/tests/testutils/fixtures.go: Funções para criar dados de teste (CreateTestUser, etc.).
  - SupplyInventory/tests/integration/helpers.go: IntegrationTestHelper com métodos para setup de routers e fixtures.
  - Use transações para isolamento e truncate para limpeza entre testes.
  - **Adicionando novas rotas:** Atualize a função `setupTestDependencies()` e crie método `SetupRouterFor{NewController}` - não é necessário modificar main.go.

  **Estrutura do IntegrationTestHelper:**
  - `setupTestDependencies(db)`: Função central que cria todas as dependências (repos, serviços) com DB de teste
  - `SetupRouterFor{Controller}(db)`: Métodos específicos que usam as dependências para configurar routers individuais
  - `SetupTestRouter(db)`: Router completo com todas as rotas para testes abrangentes
  - **Project testing patterns (recommended)**:
    - Repositories must accept a `dbadapter.DBAdapter` interface (see `SupplyInventory/Infrastructure/repositories/db/adapter.go`). This decouples business code from GORM and makes unit tests simpler.
    - Use a shared fake adapter for unit tests: `SupplyInventory/tests/testutils/fake_db_adapter.go`. It provides hook functions (`CreateFn`, `FirstByIDFn`, `WhereFirstFn`, etc.) that tests can set to control repository behavior without touching the database.
    - Integration tests should use `IntegrationTestHelper` (`SupplyInventory/tests/integration/helpers.go`) that sets up the test DB, runs migrations, truncates tables and provides router setup helpers that wire repositories using `dbadapter.NewGormAdapter(helper.DB)`.
    - Follow the test structure blocks for unit tests: `// Set`, `// Expectations` (only when using mocks/hooks), `// Actions`, `// Assertions`.
    - Branch and commit policy for test work: create focused branches named `test/<scope>` (for example `test/repositories-add`). Commit each test file only after verifying it passes locally. Use `test:` in the commit subject for test-only changes (e.g. `test(pallet): add unit tests for PalletRepository`).
    - Example dev-container commands:
      - Run unit tests: `docker exec -it go_inventory_dev /usr/local/go/bin/go test ./... -v`
      - Run integration tests (example with MySQL container IP):
        `docker exec -e TEST_DB_HOST=172.17.0.2 -e TEST_DB_PORT=3306 -e TEST_DB_USER=root -e TEST_DB_PASSWORD=root go_inventory_dev /usr/local/go/bin/go test -tags integration ./... -v`
