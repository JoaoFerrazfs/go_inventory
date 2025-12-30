# Main Rules to follow
---
### General
  - When you are chatting, always respond in brasilian portuguese, unless I explicitly ask you to respond in another language.
  - Comments in code must be written in English.
  - Write clean, readable, and maintainable code following established coding standards.
  - Use meaningful variable and function names that convey purpose.
  - Keep functions and methods small and focused on a single task.
  - Comment code where necessary to explain complex logic, but avoid over-commenting.
  - Use version control (e.g., Git) effectively with clear commit messages.
  - Conduct code reviews to ensure quality and share knowledge among team members.
  - Continuously refactor code to improve structure and performance.
  - Stay updated with the latest developments in relevant technologies and best practices.
  - Prioritize security in all aspects of development.
  - Write tests to cover critical functionality and edge cases.
### APIs
  - Use RESTful API design principles.
  - Use consistent naming conventions for endpoints (e.g., plural nouns).
  - Implement proper HTTP status codes for responses.
  - Ensure endpoints are well-documented with clear descriptions, parameters, and response formats.
  - Validate all incoming data and handle errors gracefully.
  - Implement versioning for APIs to manage changes over time.
  - Use authentication and authorization mechanisms to secure endpoints.
  - Optimize for performance, including pagination for large datasets.
  - Write unit and integration tests for API endpoints to ensure reliability.
  - Follow security best practices to protect against common vulnerabilities (e.g., SQL injection, XSS).
### Documentation
  - Maintain comprehensive documentation for codebases, including setup instructions, usage guides, and API references.
  - Use clear and concise language in documentation.
  - Keep documentation up-to-date with code changes.
  - Include examples and use cases to illustrate functionality.
  - Organize documentation logically with a table of contents for easy navigation.
  - Use diagrams and visuals where appropriate to enhance understanding.
  - Encourage contributions to documentation from all team members.
  - Review and revise documentation regularly to ensure accuracy.
  - Host documentation in this folder "/documents".
  - Each document must be organized following the same paste structure as the other project files
    - For example following files, where each related document is stored in:
      - /documents/domain/Pallet
      - /documents/domain/PalletProduct
      - /documents/domain/PalletRackEntity
      - /documents/domain/User
  - Provide troubleshooting guides for common issues.  
### Testing
  - Os testes devem ser sempre faceis de entender com a leitura facilitada, estruturado e organizado.
  - Não use siglas mas sim as palavras completas para nomes de variáveis, funções, métodos e arquivos de teste.
   - Prefira leitura simples ao invés de complexa, mesmo que a complexa seja mais "elegante".
  - Para rodar os testes Go, utilize sempre o container Docker da aplicação.
  - Comandos principais (use o binário `go` absoluto dentro do container para maior confiabilidade):
    - Testes unitários (com mocks):
      - docker exec -it go_inventory_dev /usr/local/go/bin/go test ./SupplyInventory/Application/Controllers/Auth/... -v
    - Testes de integração (com banco de dados real):
      - docker exec -e TEST_DB_HOST=172.17.0.2 -e TEST_DB_PORT=3306 -e TEST_DB_USER=root -e TEST_DB_PASSWORD=root \
        go_inventory_dev /usr/local/go/bin/go test -tags integration ./SupplyInventory/Application/Controllers/Auth/... -v
    - Todos os testes unitários:
      - docker exec -it go_inventory_dev /usr/local/go/bin/go test ./... -v
    - Todos os testes de integração:
      - docker exec -e TEST_DB_HOST=172.17.0.2 -e TEST_DB_PORT=3306 -e TEST_DB_USER=root -e TEST_DB_PASSWORD=root \
        go_inventory_dev /usr/local/go/bin/go test -tags integration ./... -v

    - Run both (unit then integration) in sequence:
      - docker exec -it go_inventory_dev /usr/local/go/bin/go test ./... -v && \
        docker exec -e TEST_DB_HOST=172.17.0.2 -e TEST_DB_PORT=3306 -e TEST_DB_USER=root -e TEST_DB_PASSWORD=root \
        go_inventory_dev /usr/local/go/bin/go test -tags integration ./... -v
  - Certifique-se de instalar as dependências de teste no ambiente do container:
    - github.com/stretchr/testify
    - github.com/davecgh/go-spew/spew
    - github.com/pmezard/go-difflib/difflib
    - github.com/stretchr/objx
  - Use mocks do testify para simular serviços e repositórios nos testes unitários.
  - Para testes de integração, use banco de dados MySQL real com containers.
  - Siga a estrutura de pastas de testes conforme o domínio do código.
  - Write unit tests for individual components and functions.
  - Implement integration tests to verify interactions between components.
    - Each API endpoint must have an integration test.
  - Use end-to-end testing to simulate real user scenarios.
  - Aim for high test coverage, focusing on critical and complex areas of the codebase.
  - Use testing frameworks and tools appropriate for the technology stack.
  - Write clear and descriptive test cases that outline expected behavior.
  - Perform load and performance testing to ensure the application can handle expected traffic.
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
    - Deve haver uma linha em branco entre cada bloco.
    - O bloco // Expectations só deve aparecer se houver expectativas de mocks; caso contrário, omita.
    - Todos os comentários devem ser em inglês.

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
### Commits
  - Use clear and descriptive commit messages that summarize the changes made.
  - Follow a consistent format for commit messages, such as:
    - feat: for new features
    - fix: for bug fixes
    - docs: for documentation updates
    - style: for code style changes (formatting, etc.)
    - refactor: for code restructuring without changing functionality
    - test: for adding or updating tests
    - chore: for maintenance tasks (build process, dependencies, etc.)
  - Keep commits small and focused on a single change or feature.
  - Avoid committing generated files or build artifacts.
  - Use branches effectively to manage features, bug fixes, and releases.
  - Regularly pull changes from the main branch to keep your branch up-to-date.
  - Review changes before committing to ensure code quality.
  - Keep small commits, but with descriptive messages.
  - Use imperative mood in commit messages (e.g., "Add feature" instead of "Added feature").
  - Reference relevant issue numbers in commit messages when applicable.
