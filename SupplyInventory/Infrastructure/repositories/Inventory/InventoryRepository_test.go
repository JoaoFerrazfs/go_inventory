package infrastructure

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	entities "go_inventory/SupplyInventory/Domain/Entities"
	"go_inventory/SupplyInventory/tests/testutils"
)

func TestInventoryRepository_List_Success(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repository := NewInventoryRepository(adapter)

	// Actions
	inventories, appErr := repository.List()

	// Assertions
	assert.Nil(t, appErr)
	assert.Empty(t, inventories)
}

func TestInventoryRepository_List_Error(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{
		FindAllFn: func(out interface{}) error {
			return errors.New("database error")
		},
	}
	repository := NewInventoryRepository(adapter)

	// Actions
	inventories, appErr := repository.List()

	// Assertions
	assert.Nil(t, inventories)
	assert.NotNil(t, appErr)
	assert.Equal(t, 500, appErr.ErrorCode())
	assert.Equal(t, "Failed to list inventories", appErr.Error())
}

func TestInventoryRepository_Exists_Success(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repository := NewInventoryRepository(adapter)

	// Actions
	exists, appErr := repository.Exists(1)

	// Assertions
	assert.Nil(t, appErr)
	assert.True(t, exists)
}

func TestInventoryRepository_Exists_NotFound(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{
		FirstByIDFn: func(out interface{}, id uint) error {
			return errors.New("not found")
		},
	}
	repository := NewInventoryRepository(adapter)

	// Actions
	exists, appErr := repository.Exists(999)

	// Assertions
	assert.NotNil(t, appErr)
	assert.False(t, exists)
	assert.Equal(t, 404, appErr.ErrorCode())
}

func TestInventoryRepository_FindById_Success(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repository := NewInventoryRepository(adapter)

	// Actions
	inventory, appErr := repository.FindById(1)

	// Assertions
	assert.Nil(t, appErr)
	assert.NotNil(t, inventory)
}

func TestInventoryRepository_FindById_NotFound(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{
		FirstByIDFn: func(out interface{}, id uint) error {
			return errors.New("not found")
		},
	}
	repository := NewInventoryRepository(adapter)

	// Actions
	inventory, appErr := repository.FindById(999)

	// Assertions
	assert.Nil(t, inventory)
	assert.NotNil(t, appErr)
	assert.Equal(t, 404, appErr.ErrorCode())
}

func TestInventoryRepository_Create_Success(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repository := NewInventoryRepository(adapter)
	inv := &entities.InventoryEntity{UserID: 1, Status: entities.InventoryStatusOpen}

	// Actions
	appErr := repository.Create(inv)

	// Assertions
	assert.Nil(t, appErr)
}

func TestInventoryRepository_Create_Error(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{
		CreateFn: func(value interface{}) error {
			return errors.New("create failed")
		},
	}
	repository := NewInventoryRepository(adapter)
	inv := &entities.InventoryEntity{UserID: 1, Status: entities.InventoryStatusOpen}

	// Actions
	appErr := repository.Create(inv)

	// Assertions
	assert.NotNil(t, appErr)
	assert.Equal(t, 500, appErr.ErrorCode())
}

func TestInventoryRepository_Update_Success(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repository := NewInventoryRepository(adapter)
	inv := &entities.InventoryEntity{ID: 1, UserID: 1, Status: entities.InventoryStatusClosed}

	// Actions
	appErr := repository.Update(inv)

	// Assertions
	assert.Nil(t, appErr)
}

func TestInventoryRepository_Update_Error(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{
		SaveFn: func(value interface{}) error {
			return errors.New("update failed")
		},
	}
	repository := NewInventoryRepository(adapter)
	inv := &entities.InventoryEntity{ID: 1, UserID: 1, Status: entities.InventoryStatusClosed}

	// Actions
	appErr := repository.Update(inv)

	// Assertions
	assert.NotNil(t, appErr)
	assert.Equal(t, 500, appErr.ErrorCode())
}
