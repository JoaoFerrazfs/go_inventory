# Testing Standards

## Overview
This document outlines the testing standards and practices for the go_inventory project. We maintain both unit tests (with fakes/mocks) and integration tests (with a real database) to ensure code quality and reliability.

This project follows a small set of patterns to keep unit tests fast and maintainable while keeping integration tests realistic:

- Use a `DBAdapter` interface so repositories are not coupled to GORM types directly.
- Use a shared `FakeDBAdapter` in unit tests with per-test hooks to control returned data and errors.
- Use an `IntegrationTestHelper` that provides a real `*gorm.DB`, runs migrations, truncates tables, and exposes router setup helpers.

## Test Structure

### Unit Tests
- **Location**: Same directory as the code being tested
- **Naming**: `{FileName}_test.go`
- **Purpose**: Test individual functions and methods in isolation
- **Mocking**: Use `testify` mocks or the shared `FakeDBAdapter` for repositories

Recommended repository unit-test pattern:

- Repositories must depend on `dbadapter.DBAdapter` (not `*gorm.DB`).
- Use `SupplyInventory/tests/testutils.FakeDBAdapter` and set only the hooks needed for the test (e.g. `FirstByIDFn`, `WhereFirstFn`, `CreateFn`).
- Follow structure:

```go
func TestSomething_Scenario(t *testing.T) {
    // Set
    fake := &testutils.FakeDBAdapter{}
    repo := repositories.NewSomethingRepository(fake)

    // Expectations
    fake.FirstByIDFn = func(dest interface{}, id uint) error {
        // populate dest
        return nil
    }

    // Actions
    got, err := repo.FindByID(1)

    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, expected, got)
}
```

### Integration Tests
- **Location**: Same directory as the code being tested
- **Naming**: `{FileName}_integration_test.go`
- **Build Tags**: `//go:build integration`
- **Purpose**: Test API endpoints and database interactions
- **Database**: Real MySQL database with transactions

Integration test pattern:

- Mark files with `//go:build integration` and use `IntegrationTestHelper` (`SupplyInventory/tests/integration/helpers.go`).
- Each test should call `h.TruncateTables(h.DB)` before executing and run test logic inside `h.DB.Transaction(...)` for isolation.
- Use `SetupRouterFor{Controller}` to obtain a router wired with `dbadapter.NewGormAdapter(h.DB)`.

Example:

```go
//go:build integration

func TestIntegration_CreateUser(t *testing.T) {
    h := integration.NewIntegrationTestHelper()
    h.TruncateTables(h.DB)

    h.DB.Transaction(func(tx *gorm.DB) error {
        r := h.SetupRouterForUser(tx)
        // build request and call r.ServeHTTP
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

### DBAdapter and FakeDBAdapter

- `SupplyInventory/Infrastructure/repositories/db/adapter.go` defines the `DBAdapter` interface and a small gorm-backed implementation `NewGormAdapter(db *gorm.DB)`.
- `SupplyInventory/tests/testutils/fake_db_adapter.go` contains `FakeDBAdapter` with hook fields (e.g. `CreateFn`, `FirstByIDFn`, `WhereFirstFn`, `FindAllFn`, `DeleteByIDFn`, `SaveFn`, `PreloadFindFn`, `AppendAssociationFn`). Unit tests should set only the hooks they need.

Benefits:
- Tests are simpler (no sqlmock / complicated gorm expectations).
- Tests run fast and avoid touching the network for unit tests.

Example hook usage:

```go
fake := &testutils.FakeDBAdapter{}
fake.WhereFirstFn = func(dest interface{}, query string, args ...interface{}) error {
    // fill dest expected struct
    return nil
}
```

## Best Practices

### General
- Tests should be fast and reliable
- Use descriptive test names: `Test{Function}_{Scenario}`
- Keep tests focused on one behavior
- Avoid test interdependencies

Additional project rules to follow for all tests:

- Branching & commits for tests: create a short-lived branch `test/<what>` for a group of tests (example: `test/repositories-add`). Commit each new test only after you run it and it passes locally. Use `test:` prefix in commit messages when adding tests (e.g. `test(repo): add unit tests for PalletRepository`).
- Run unit tests frequently (`go test ./...`) and run integration tests (`go test -tags integration ./...`) before creating a PR.
- Integration tests should run in CI as a separate job with a real MySQL container and the `inventory_test` database.

## Dependencies
Ensure these packages are installed in the test container:
- `github.com/stretchr/testify`
- `github.com/davecgh/go-spew/spew`
- `github.com/pmezard/go-difflib/difflib`
- `github.com/stretchr/objx`

## Quick Commands

Run unit tests inside the dev container:

```bash
docker exec -it go_inventory_dev /usr/local/go/bin/go test ./... -v
```

Run integration tests inside the dev container (example using MySQL container IP):

```bash
docker exec -e TEST_DB_HOST=172.17.0.2 -e TEST_DB_PORT=3306 -e TEST_DB_USER=root -e TEST_DB_PASSWORD=root go_inventory_dev /usr/local/go/bin/go test -tags integration ./... -v
```

## Coverage Goals
- Unit tests: Focus on complex logic
- Integration tests: 80%+ coverage on API endpoints
- Combined coverage: High confidence in code reliability