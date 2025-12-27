# Main Rules to follow
---
### General
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
  - Para rodar os testes Go, utilize sempre o container Docker da aplicação.
  - Comandos principais:
    - Testes unitários:
      - docker exec -it go_inventory_dev go test ./tests/unit/SupplyInventory/Application/Controllers/...
    - Testes de integração:
      - docker exec -it go_inventory_dev go test ./tests/integration/SupplyInventory/Application/Controllers/...
    - Todos os testes:
      - docker exec -it go_inventory_dev go test ./...
  - Certifique-se de instalar as dependências de teste no ambiente do container:
    - github.com/stretchr/testify
    - github.com/davecgh/go-spew/spew
    - github.com/pmezard/go-difflib/difflib
    - github.com/stretchr/objx
  - Use mocks do testify para simular serviços e repositórios nos testes.
  - Siga a estrutura de pastas de testes conforme o domínio do código.
  - Write unit tests for individual components and functions.
  - Implement integration tests to verify interactions between components.
    - Each api need to have a integration test.
  - Use end-to-end testing to simulate real user scenarios.
  - Aim for high test coverage, focusing on critical and complex areas of the codebase.
  - Use testing frameworks and tools appropriate for the technology stack.
  - Write clear and descriptive test cases that outline expected behavior.
  - Perform load and performance testing to ensure the application can handle expected traffic.
  - Todos os arquivos de teste Go (_test.go) devem ficar na mesma pasta do arquivo que testam.
    - Exemplo: se existe /SupplyInventory/Application/Controllers/AuthController.go, o teste deve ser /SupplyInventory/Application/Controllers/AuthController_test.go.
    - Isso garante que o coverage funcione corretamente e segue o padrão da comunidade Go.

  - **Organização dos blocos dos testes:**
    - Todos os testes devem ser organizados em blocos separados e comentados, seguindo a ordem:
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
