package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go_inventory/SupplyInventory/tests/testutils"
)

func TestPalletRack_Create_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repo := NewPalletRackRepository(adapter)

	// Actions
	rack, err := repo.Create("Rack1", "Loc", 10)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, rack)
}

func TestPalletRack_ListRacks_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repo := NewPalletRackRepository(adapter)

	// Actions
	racks, err := repo.ListRacks()

	// Assertions
	assert.NoError(t, err)
	// fakeAdapter does not populate data; expect empty slice
	assert.Equal(t, 0, len(racks))
}

func TestPalletRack_FindPalletById_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repo := NewPalletRackRepository(adapter)

	// Actions
	_, appErr := repo.FindPalletById(1)

	// Assertions
	assert.Nil(t, appErr)
}

func TestPalletRack_DeleteRack_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repo := NewPalletRackRepository(adapter)

	// Actions
	ok, appErr := repo.DeleteRack(1)

	// Assertions
	assert.Nil(t, appErr)
	assert.True(t, ok)
}
