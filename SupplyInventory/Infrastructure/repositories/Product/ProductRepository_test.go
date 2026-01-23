package infrastructure_test

import (
	"testing"

	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	repositories "go_inventory/SupplyInventory/Infrastructure/repositories/Product"

	"github.com/stretchr/testify/assert"

	testutils "go_inventory/SupplyInventory/tests/testutils"
)

func TestProductRepository_Create_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repo := repositories.NewProductRepository(adapter)

	// Expectations
	adapter.SaveFn = func(value interface{}) error {
		productPtr := value.(**entities.ProductEntity)
		(*productPtr).ID = 1
		return nil
	}

	// Actions
	product, _ := repo.Create("123456789", "Test Product")

	// Assertions
	assert.NotNil(t, product)
	assert.Equal(t, uint(1), product.ID)
	assert.Equal(t, "123456789", product.EAN)
	assert.Equal(t, "Test Product", product.Name)
}

func TestProductRepository_FindByEAN_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repo := repositories.NewProductRepository(adapter)

	// Expectations
	adapter.WhereFirstFn = func(out interface{}, query string, args ...interface{}) error {
		if query == "ean = ?" && args[0] == "123456789" {
			*out.(*entities.ProductEntity) = entities.ProductEntity{ID: 1, EAN: "123456789", Name: "Test Product"}
		}
		return nil
	}

	// Actions
	product, _ := repo.FindByEAN("123456789")

	// Assertions
	assert.NotNil(t, product)
	assert.Equal(t, "123456789", product.EAN)
	assert.Equal(t, "Test Product", product.Name)
}

func TestProductRepository_FindByEAN_NotFound_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repo := repositories.NewProductRepository(adapter)

	// Expectations
	adapter.WhereFirstFn = func(out interface{}, query string, args ...interface{}) error {
		return errors.NewAppError("record not found", 404)
	}

	// Actions
	product, err := repo.FindByEAN("nonexistent")

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, "record not found", err.Message)
	assert.Nil(t, product)
}

func TestProductRepository_Delete_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repo := repositories.NewProductRepository(adapter)

	// Expectations
	adapter.WhereFirstFn = func(out interface{}, query string, args ...interface{}) error {
		if query == "ean = ?" && args[0] == "123456789" {
			*out.(*entities.ProductEntity) = entities.ProductEntity{ID: 1, EAN: "123456789", Name: "Test Product"}
		}
		return nil
	}
	adapter.DeleteByIDFn = func(model interface{}, id uint) (int64, error) {
		return 1, nil
	}

	// Actions
	deleted, err := repo.Delete("123456789")

	// Assertions
	assert.Nil(t, err)
	assert.True(t, deleted)
}

func TestProductRepository_Delete_NotFound_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repo := repositories.NewProductRepository(adapter)

	// Expectations
	adapter.WhereFirstFn = func(out interface{}, query string, args ...interface{}) error {
		return errors.NewAppError("record not found", 404)
	}

	// Actions
	deleted, err := repo.Delete("nonexistent")

	// Assertions
	assert.Nil(t, err)
	assert.False(t, deleted)
}
