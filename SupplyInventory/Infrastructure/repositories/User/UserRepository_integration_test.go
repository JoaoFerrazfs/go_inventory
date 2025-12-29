//go:build integration
// +build integration

package infrastructure_test

import (
	"testing"

	entities "go_inventory/SupplyInventory/Domain/Entities"
	userInfra "go_inventory/SupplyInventory/Infrastructure/repositories/User"
	dbadapter "go_inventory/SupplyInventory/Infrastructure/repositories/db"
	integration "go_inventory/SupplyInventory/tests/integration"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create_Integration(t *testing.T) {
    // Set
    helper := integration.NewIntegrationTestHelper()
    defer helper.Stop()
    helper.TruncateTables(helper.DB)
    repo := userInfra.NewUserRepository(dbadapter.NewGormAdapter(helper.DB))

    u := &entities.UserEntity{Name: "IntUser", Email: "int@example.com", Password: "pwd"}

    // Actions
    err := repo.Create(u)

    // Assertions
    assert.Nil(t, err)
    assert.NotZero(t, u.ID)
}

func TestUserRepository_FindByEmail_Integration(t *testing.T) {
    // Set
    helper := integration.NewIntegrationTestHelper()
    defer helper.Stop()
    helper.TruncateTables(helper.DB)
    // create fixture
    u := &entities.UserEntity{Name: "IntUser2", Email: "findme@example.com", Password: "pwd"}
    helper.DB.Create(u)

    repo := userInfra.NewUserRepository(dbadapter.NewGormAdapter(helper.DB))

    // Actions
    got, err := repo.FindByEmail("findme@example.com")

    // Assertions
    assert.Nil(t, err)
    assert.NotNil(t, got)
    assert.Equal(t, "findme@example.com", got.Email)
}
