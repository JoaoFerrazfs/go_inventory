//go:build integration
// +build integration

package infrastructure_test

import (
	"testing"

	palletRackInfra "go_inventory/SupplyInventory/Infrastructure/repositories/PalletRack"
	dbadapter "go_inventory/SupplyInventory/Infrastructure/repositories/db"
	integration "go_inventory/SupplyInventory/tests/integration"

	"github.com/stretchr/testify/assert"
)

func TestPalletRack_Create_Integration(t *testing.T) {
	// Set
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)
	repo := palletRackInfra.NewPalletRackRepository(dbadapter.NewGormAdapter(helper.DB))

	// Actions
	rack, err := repo.Create("RackIntegration", "Loc", 10)

	// Assertions
	assert.Nil(t, err)
	assert.NotNil(t, rack)
	assert.Equal(t, "RackIntegration", rack.Name)
}

func TestPalletRack_ListRacks_Integration(t *testing.T) {
	// Set
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)
	repo := palletRackInfra.NewPalletRackRepository(dbadapter.NewGormAdapter(helper.DB))
	helper.CreateTestPalletRack(helper.DB, "R1", "L", 5)
	helper.CreateTestPalletRack(helper.DB, "R2", "L2", 6)

	// Actions
	racks, err := repo.ListRacks()

	// Assertions
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(racks), 2)
}

func TestPalletRack_FindPalletById_Integration(t *testing.T) {
	// Set
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)
	repo := palletRackInfra.NewPalletRackRepository(dbadapter.NewGormAdapter(helper.DB))
	rack := helper.CreateTestPalletRack(helper.DB, "RackFind", "Loc", 5)
	// create a pallet to ensure preload works
	helper.CreateTestPallet(helper.DB, "PalletOnRack", rack.ID)

	// Actions
	got, appErr := repo.FindPalletById(rack.ID)

	// Assertions
	assert.Nil(t, appErr)
	assert.NotNil(t, got)
	assert.Equal(t, rack.ID, got.ID)
}

func TestPalletRack_DeleteRack_Integration(t *testing.T) {
	// Set
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)
	repo := palletRackInfra.NewPalletRackRepository(dbadapter.NewGormAdapter(helper.DB))

	// Case 1: delete empty rack (success)
	empty := helper.CreateTestPalletRack(helper.DB, "RackDel", "Loc", 5)
	ok, err := repo.DeleteRack(empty.ID)
	assert.Nil(t, err)
	assert.True(t, ok)

	// Case 2: delete rack that has pallets -> expect error
	r2 := helper.CreateTestPalletRack(helper.DB, "RackDelHas", "Loc", 5)
	helper.CreateTestPallet(helper.DB, "P1", r2.ID)
	ok2, err2 := repo.DeleteRack(r2.ID)
	assert.NotNil(t, err2)
	assert.False(t, ok2)
}
