//go:build integration
// +build integration

package infrastructure

import (
	"testing"

	entities "go_inventory/SupplyInventory/Domain/Entities"
	dbadapter "go_inventory/SupplyInventory/Infrastructure/repositories/db"
	integration "go_inventory/SupplyInventory/tests/integration"

	"github.com/stretchr/testify/assert"
)

func TestPalletRepository_AddAndGetIntegration(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	defer h.Stop()

	h.TruncateTables(h.DB)

	// create pallet rack fixture
	rack := h.CreateTestPalletRack(h.DB, "Rack1", "Loc", 10)

	// add pallet using repository bound to the test DB
	repo := NewPalletRepository(dbadapter.NewGormAdapter(h.DB))
	pallet, appErr := repo.AddSupply("IntegrationPallet", rack.ID)
	if appErr != nil {
		t.Fatalf("failed to add supply: %v", appErr)
	}

	got, appErr := repo.GetSupplyById(pallet.ID)
	if appErr != nil {
		t.Fatalf("failed to get supply: %v", appErr)
	}

	assert.Equal(t, "IntegrationPallet", got.Name)
	// create product and associate
	product := entities.PalletizedProductEntity{ProductCode: "P1", PalletID: got.ID}
	updated, appErr := repo.AddProductsToPallet(product)
	if appErr != nil {
		t.Fatalf("failed to add product: %v", appErr)
	}

	assert.NotNil(t, updated.PalletizedProduct)
}
