//go:build integration
// +build integration

package infrastructure_test

import (
	"testing"

	entities "go_inventory/SupplyInventory/Domain/Entities"
	palletInfra "go_inventory/SupplyInventory/Infrastructure/repositories/Pallet"
	dbadapter "go_inventory/SupplyInventory/Infrastructure/repositories/db"
	integration "go_inventory/SupplyInventory/tests/integration"

	"github.com/stretchr/testify/assert"
)

func TestPalletRepository_AddAndGetIntegration(t *testing.T) {
	t.Skip("deprecated combined test: replaced by focused tests")
}

func TestPalletRepository_AddSupply_Integration(t *testing.T) {
	// Set
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)
	rack := helper.CreateTestPalletRack(helper.DB, "Rack_Add", "Loc", 5)
	repo := palletInfra.NewPalletRepository(dbadapter.NewGormAdapter(helper.DB))

	// Actions
	created, err := repo.AddSupply("Pallet_Add", rack.ID)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, "Pallet_Add", created.Name)
	assert.NotZero(t, created.ID)
}

func TestPalletRepository_GetSupplyById_Integration(t *testing.T) {
	// Set
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)
	rack := helper.CreateTestPalletRack(helper.DB, "Rack_Get", "Loc", 5)
	repo := palletInfra.NewPalletRepository(dbadapter.NewGormAdapter(helper.DB))
	created, err := repo.AddSupply("Pallet_Get", rack.ID)
	if err != nil {
		t.Fatalf("setup AddSupply failed: %v", err)
	}

	// Actions
	got, getErr := repo.GetSupplyById(created.ID)

	// Assertions
	assert.Nil(t, getErr)
	assert.Equal(t, "Pallet_Get", got.Name)
}

func TestPalletRepository_AddProductsToPallet_Integration(t *testing.T) {
	// Set
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)
	rack := helper.CreateTestPalletRack(helper.DB, "Rack_AddProd", "Loc", 5)
	repo := palletInfra.NewPalletRepository(dbadapter.NewGormAdapter(helper.DB))
	created, err := repo.AddSupply("Pallet_AddProd", rack.ID)
	if err != nil {
		t.Fatalf("setup AddSupply failed: %v", err)
	}

	// Actions
	product := entities.PalletizedProductEntity{EAN: 100, Quantity: 1, PalletID: created.ID}
	updated, addErr := repo.AddProductsToPallet(product)

	// Assertions
	assert.Nil(t, addErr)
	assert.NotNil(t, updated.PalletizedProduct)
	assert.Greater(t, len(updated.PalletizedProduct), 0)
}
