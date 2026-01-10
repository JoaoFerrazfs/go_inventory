package infrastructure_test

import (
	"fmt"
	"testing"

	entities "go_inventory/SupplyInventory/Domain/Entities"
	inventoryInfra "go_inventory/SupplyInventory/Infrastructure/repositories/Inventory"
	"go_inventory/SupplyInventory/tests/testutils"

	"github.com/stretchr/testify/assert"
)

func TestInventoryRepository_Exists_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	adapter.FirstByIDFn = func(out interface{}, id uint) error {
		if id == 1 {
			return nil
		}
		return fmt.Errorf("not found")
	}
	repo := inventoryInfra.NewInventoryRepository(adapter)

	// Actions & Assertions
	exists, err := repo.Exists(1)
	assert.True(t, exists)
	assert.Nil(t, err)

	exists, err = repo.Exists(2)
	assert.False(t, exists)
	assert.NotNil(t, err)
	assert.Equal(t, 404, err.Code)
}

func TestInventoryRepository_Create_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repo := inventoryInfra.NewInventoryRepository(adapter)
	inv := &entities.InventoryEntity{Name: "New Inventory", UserID: 1}

	// Actions
	err := repo.Create(inv)

	// Assertions
	assert.Nil(t, err)
}

func TestInventoryRepository_FindById_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	adapter.PreloadFindFn = func(out interface{}, preload string, id ...uint) error {
		if id[0] == 1 {
			inv := out.(*entities.InventoryEntity)
			inv.ID = 1
			inv.Name = "Found"
			return nil
		}
		return fmt.Errorf("not found")
	}
	repo := inventoryInfra.NewInventoryRepository(adapter)

	// Actions & Assertions
	inv, err := repo.FindById(1)
	assert.Nil(t, err)
	assert.NotNil(t, inv)
	assert.Equal(t, uint(1), inv.ID)

	inv, err = repo.FindById(2)
	assert.NotNil(t, err)
	assert.Nil(t, inv)
	assert.Equal(t, 404, err.Code)
}

func TestInventoryRepository_Update_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repo := inventoryInfra.NewInventoryRepository(adapter)
	inv := &entities.InventoryEntity{ID: 1, Name: "To Update"}

	// Actions
	err := repo.Update(inv)

	// Assertions
	assert.Nil(t, err)
}
