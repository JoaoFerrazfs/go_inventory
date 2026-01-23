//go:build integration
// +build integration

package infrastructure_test

import (
	"testing"

	productInfra "go_inventory/SupplyInventory/Infrastructure/repositories/Product"
	dbadapter "go_inventory/SupplyInventory/Infrastructure/repositories/db"
	integration "go_inventory/SupplyInventory/tests/integration"

	"github.com/stretchr/testify/assert"
)

func TestProductRepository_Create_Integration(t *testing.T) {
	// Set
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)
	repo := productInfra.NewProductRepository(dbadapter.NewGormAdapter(helper.DB))

	// Actions
	product, err := repo.Create("1234567890123", "Integration Product")

	// Assertions
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.NotZero(t, product.ID)
	assert.Equal(t, "1234567890123", product.EAN)
	assert.Equal(t, "Integration Product", product.Name)
}

func TestProductRepository_FindByEAN_Integration(t *testing.T) {
	// Set
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)
	repo := productInfra.NewProductRepository(dbadapter.NewGormAdapter(helper.DB))

	// Create a product first
	createdProduct := helper.CreateTestProduct(helper.DB, "1234567890123", "Find Test Product")

	// Actions
	foundProduct, err := repo.FindByEAN("1234567890123")

	// Assertions
	assert.Nil(t, err)
	assert.NotNil(t, foundProduct)
	assert.Equal(t, createdProduct.ID, foundProduct.ID)
	assert.Equal(t, "1234567890123", foundProduct.EAN)
	assert.Equal(t, "Find Test Product", foundProduct.Name)
}

func TestProductRepository_FindByEAN_NotFound_Integration(t *testing.T) {
	// Set
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)
	repo := productInfra.NewProductRepository(dbadapter.NewGormAdapter(helper.DB))

	// Actions
	product, err := repo.FindByEAN("nonexistent")

	// Assertions
	assert.NotNil(t, err)
	assert.Nil(t, product)
}

func TestProductRepository_Delete_Integration(t *testing.T) {
	// Set
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)
	repo := productInfra.NewProductRepository(dbadapter.NewGormAdapter(helper.DB))

	// Create a product first
	helper.CreateTestProduct(helper.DB, "1234567890123", "Delete Test Product")

	// Actions
	deleted, err := repo.Delete("1234567890123")

	// Assertions
	assert.Nil(t, err)
	assert.True(t, deleted)

	// Verify it's deleted
	product, err := repo.FindByEAN("1234567890123")
	assert.NotNil(t, err)
	assert.Nil(t, product)
}

func TestProductRepository_Delete_NotFound_Integration(t *testing.T) {
	// Set
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)
	repo := productInfra.NewProductRepository(dbadapter.NewGormAdapter(helper.DB))

	// Actions
	deleted, err := repo.Delete("nonexistent")

	// Assertions
	assert.Nil(t, err)
	assert.False(t, deleted)
}
