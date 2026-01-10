//go:build integration
// +build integration

package infrastructure_test

import (
	"testing"
	"time"

	entities "go_inventory/SupplyInventory/Domain/Entities"
	inventoryInfra "go_inventory/SupplyInventory/Infrastructure/repositories/Inventory"
	dbadapter "go_inventory/SupplyInventory/Infrastructure/repositories/db"
	integration "go_inventory/SupplyInventory/tests/integration"

	"github.com/stretchr/testify/assert"
)

func TestInventoryRepository_Create_Integration(t *testing.T) {
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)

	user := helper.CreateTestUser(helper.DB, "Repo User", "repo@example.com", "password")
	repo := inventoryInfra.NewInventoryRepository(dbadapter.NewGormAdapter(helper.DB))

	inv := &entities.InventoryEntity{
		Name:      "Integration Inventory",
		UserID:    user.ID,
		Status:    entities.InventoryStatusOpen,
		StartedAt: time.Now(),
	}

	// Action
	err := repo.Create(inv)

	// Assertions
	assert.Nil(t, err)
	assert.NotZero(t, inv.ID)
}

func TestInventoryRepository_FindById_Integration(t *testing.T) {
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)

	invFixture := helper.CreateTestInventory(helper.DB)
	repo := inventoryInfra.NewInventoryRepository(dbadapter.NewGormAdapter(helper.DB))

	// Action
	inv, err := repo.FindById(invFixture.ID)

	// Assertions
	assert.Nil(t, err)
	assert.NotNil(t, inv)
	assert.Equal(t, invFixture.Name, inv.Name)
	assert.Equal(t, invFixture.UserID, inv.UserID)
	assert.NotNil(t, inv.User)
}

func TestInventoryRepository_Exists_Integration(t *testing.T) {
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)

	invFixture := helper.CreateTestInventory(helper.DB)
	repo := inventoryInfra.NewInventoryRepository(dbadapter.NewGormAdapter(helper.DB))

	// Actions
	exists, err := repo.Exists(invFixture.ID)
	assert.True(t, exists)
	assert.Nil(t, err)

	exists, err = repo.Exists(999)
	assert.False(t, exists)
	assert.NotNil(t, err)
}

func TestInventoryRepository_Where_Integration(t *testing.T) {
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)

	invFixture := helper.CreateTestInventory(helper.DB)
	repo := inventoryInfra.NewInventoryRepository(dbadapter.NewGormAdapter(helper.DB))

	// Action
	inventories, err := repo.Where(map[string]any{"name": invFixture.Name})

	// Assertions
	assert.Nil(t, err)
	assert.Len(t, inventories, 1)
	assert.Equal(t, invFixture.ID, inventories[0].ID)
	assert.NotNil(t, inventories[0].User)
}

func TestInventoryRepository_Update_Integration(t *testing.T) {
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)

	invFixture := helper.CreateTestInventory(helper.DB)
	repo := inventoryInfra.NewInventoryRepository(dbadapter.NewGormAdapter(helper.DB))

	invFixture.Name = "Updated Name"

	// Action
	err := repo.Update(invFixture)

	// Assertions
	assert.Nil(t, err)

	var inv entities.InventoryEntity
	helper.DB.First(&inv, invFixture.ID)
	assert.Equal(t, "Updated Name", inv.Name)
}
