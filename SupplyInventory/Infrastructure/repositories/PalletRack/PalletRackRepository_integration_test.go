//go:build integration
// +build integration

package infrastructure_test

import (
	"testing"

	palletRackInfra "go_inventory/SupplyInventory/Infrastructure/repositories/PalletRack"
	dbadapter "go_inventory/SupplyInventory/Infrastructure/repositories/db"
	integration "go_inventory/SupplyInventory/tests/integration"
	testutils "go_inventory/SupplyInventory/tests/testutils"

	entities "go_inventory/SupplyInventory/Domain/Entities"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPalletRack_Create_Integration(t *testing.T) {
	// Set
	helper := integration.NewIntegrationTestHelper()
	defer helper.Stop()
	helper.TruncateTables(helper.DB)
	repo := palletRackInfra.NewPalletRackRepository(dbadapter.NewGormAdapter(helper.DB))

	// Actions
	// ensure inventory exists
	inv := testutils.CreateTestInventory(helper.DB)
	rack, err := repo.Create("RackIntegration", "Loc", 10, inv.ID)

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

	// Ensure racks are created successfully
	rack1 := helper.CreateTestPalletRack(helper.DB, "R1", "L", 5)
	assert.NotNil(t, rack1, "Failed to create rack R1")
	rack2 := helper.CreateTestPalletRack(helper.DB, "R2", "L2", 6)
	assert.NotNil(t, rack2, "Failed to create rack R2")

	// Actions
	racks, err := repo.ListRacks()

	// Assertions
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(racks), 2, "Expected at least 2 racks, got %d", len(racks))
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

	// Ensure rack is created successfully
	rack := helper.CreateTestPalletRack(helper.DB, "RackDel", "Loc", 5)
	assert.NotNil(t, rack, "Failed to create rack for deletion test")

	// Actions
	deleted, appErr := repo.DeleteRack(rack.ID)

	// Assertions
	assert.True(t, deleted, "Expected rack to be deleted, but it was not")
	assert.Nil(t, appErr, "Unexpected application error during deletion")

	// Verify deletion
	var deletedRack entities.PalletRackEntity
	result := helper.DB.First(&deletedRack, rack.ID)
	assert.Error(t, result.Error, "Expected error when fetching deleted rack, got none")
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error, "Expected record not found error, got: %v", result.Error)
}
