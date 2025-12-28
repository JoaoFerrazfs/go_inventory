# Testing Standards

## Overview
This document outlines the testing standards and practices for the go_inventory project. We maintain both unit tests (with mocks) and integration tests (with real database) to ensure code quality and reliability.

## Test Structure

### Unit Tests
- **Location**: Same directory as the code being tested
- **Naming**: `{FileName}_test.go`
- **Purpose**: Test individual functions and methods in isolation
- **Mocking**: Use testify mocks for external dependencies

### Integration Tests
- **Location**: Same directory as the code being tested
- **Naming**: `{FileName}_integration_test.go`
- **Build Tags**: `//go:build integration`
- **Purpose**: Test API endpoints and database interactions
- **Database**: Real MySQL database with transactions

## Running Tests

### Commands
```bash
# Unit tests
docker exec -it go_inventory_dev go test ./SupplyInventory/Application/Controllers/Auth/...

# Integration tests
docker exec -it go_inventory_dev go test -tags integration ./SupplyInventory/Application/Controllers/Auth/...

# All unit tests
docker exec -it go_inventory_dev go test ./...

# All integration tests
docker exec -it go_inventory_dev go test -tags integration ./...
```

## Test Organization

### Unit Test Structure
```go
func TestFunctionName_Scenario(t *testing.T) {
    // Set
    // Setup mocks and test data

    // Expectations
    // Define mock expectations

    // Actions
    // Execute the function

    // Assertions
    // Verify results
}
```

### Integration Test Structure
```go
//go:build integration

func TestIntegration_FunctionName(t *testing.T) {
    h := integration.NewIntegrationTestHelper()
    h.TruncateTables(h.DB)
    h.DB.Transaction(func(tx *gorm.DB) error {
        // Set
        // Setup test data and router

        // Actions
        // Make HTTP request

        // Assertions
        // Verify response
        return nil
    })
}
```

## Test Utilities

### IntegrationTestHelper
Located in `SupplyInventory/tests/integration/helpers.go`

**Purpose**: Provides consistent router setup for integration tests with proper dependency injection.

**Key Methods:**
- `NewIntegrationTestHelper()`: Creates test helper with DB connection
- `TruncateTables(db *gorm.DB)`: Cleans all test tables
- `SetupRouterFor{Controller}(db *gorm.DB)`: Creates router with test DB dependencies
- `CreateTest{Entity}(db *gorm.DB, ...)`: Creates test fixtures
- `SetupTestRouter(db *gorm.DB)`: Creates complete router with all routes (for comprehensive testing)

**Adding New Routes:**
1. Add dependencies to `setupTestDependencies()` function
2. Create `SetupRouterFor{NewController}()` method using the helper
3. Update `SetupTestRouter()` if needed for full router testing
4. **No changes needed to main.go** - dependencies are automatically handled

**Important**: When adding new controllers/routes, update the `setupTestDependencies()` function and create corresponding `SetupRouterFor{Controller}` method to maintain test consistency.

### Test Database
- **Name**: `inventory_test`
- **Setup**: Automatic via `testutils.SetupTestDB()`
- **Migration**: AutoMigrate on test DB
- **Isolation**: Each test uses transactions that rollback

### Fixtures
Located in `SupplyInventory/tests/testutils/fixtures.go`

**Available Fixtures:**
- `CreateTestUser(db, name, email, password)`
- `CreateTestPallet(db, name, palletRackID)`
- `CreateTestPalletRack(db, name, location, totalCapacity)`

## Best Practices

### General
- Tests should be fast and reliable
- Use descriptive test names: `Test{Function}_{Scenario}`
- Keep tests focused on one behavior
- Avoid test interdependencies

### Unit Tests
- Mock all external dependencies
- Test edge cases and error conditions
- Verify mock expectations with `AssertExpectations(t)`

### Integration Tests
- Test complete API workflows
- Use real database constraints and relationships
- Clean up data between tests
- Test both success and failure scenarios

### CI/CD
- Unit tests run on every push
- Integration tests run in separate job with MySQL container
- Coverage reports generated for integration tests

## Dependencies
Ensure these packages are installed in the test container:
- `github.com/stretchr/testify`
- `github.com/davecgh/go-spew/spew`
- `github.com/pmezard/go-difflib/difflib`
- `github.com/stretchr/objx`

## Coverage Goals
- Unit tests: Focus on complex logic
- Integration tests: 80%+ coverage on API endpoints
- Combined coverage: High confidence in code reliability