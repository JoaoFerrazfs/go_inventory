//go:build integration
// +build integration

package infrastructure_test

import (
	"testing"

	entities "go_inventory/SupplyInventory/Domain/Entities"
	palletizedProductInfra "go_inventory/SupplyInventory/Infrastructure/repositories/PalletizedProduct"
	dbadapter "go_inventory/SupplyInventory/Infrastructure/repositories/db"
	integration "go_inventory/SupplyInventory/tests/integration"

	"github.com/stretchr/testify/assert"
)

func TestPalletizedProduct_AddProductsToPallet_Integration(t *testing.T) {
	// Set
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)
	rack := helper.CreateTestPalletRack(helper.DB, "RackPP_Add", "Loc", 5)
	pallet := helper.CreateTestPallet(helper.DB, "PalletPP", rack.ID)
	repo := palletizedProductInfra.NewPalletizedProductRepository(dbadapter.NewGormAdapter(helper.DB), helper.PalletRepo)

	// Actions
	p := entities.PalletizedProductEntity{EAN: 555, Quantity: 2, PalletID: pallet.ID}
	ok, err := repo.AddProductsToPallet(p)

	// Assertions
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestPalletizedProduct_DeleteProductsFromPallet_Integration(t *testing.T) {
	// Set
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)
	rack := helper.CreateTestPalletRack(helper.DB, "RackPP_Del", "Loc", 5)
	pallet := helper.CreateTestPallet(helper.DB, "PalletPPDel", rack.ID)
	repo := palletizedProductInfra.NewPalletizedProductRepository(dbadapter.NewGormAdapter(helper.DB), helper.PalletRepo)
	// add product first
	prod := entities.PalletizedProductEntity{EAN: 777, Quantity: 1, PalletID: pallet.ID}
	_, addErr := repo.AddProductsToPallet(prod)
	if addErr != nil {
		t.Fatalf("setup AddProductsToPallet failed: %v", addErr)
	}

	// Actions
	ok, delErr := repo.DeleteProductsFromPallet(pallet.ID, 777)

	// Assertions
	assert.Nil(t, delErr)
	assert.True(t, ok)
}
