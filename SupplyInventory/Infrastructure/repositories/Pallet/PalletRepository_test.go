//go:build unit

package infrastructure_test

import (
	"testing"

	repository "go_inventory/SupplyInventory/Infrastructure/repositories/Pallet"

	"github.com/stretchr/testify/assert"

	entities "go_inventory/SupplyInventory/Domain/Entities"
	"go_inventory/SupplyInventory/tests/testutils"
)

func TestPalletRepository_Create_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repository := repository.NewPalletRepository(adapter)

	// Expectations
	// fakeAdapter returns success by default

	// Actions
	palletEntity := &entities.PalletEntity{Name: "TestPallet", PalletRackID: 1}
	err := repository.Create(palletEntity)

	// Assertions
	assert.Nil(t, err)
}

func TestPalletRepository_FindByID_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repository := repository.NewPalletRepository(adapter)

	// Expectations
	// fakeAdapter returns success (nil) from FirstByID

	// Actions
	_, err := repository.FindByID(1)

	// Assertions
	assert.Nil(t, err)
}

func TestPalletRepository_List_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repository := repository.NewPalletRepository(adapter)

	// Actions
	_, err := repository.List()

	// Assertions
	assert.Nil(t, err)
}

func TestPalletRepository_DeleteByID_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repository := repository.NewPalletRepository(adapter)

	// Actions
	err := repository.DeleteByID(1)

	// Assertions
	assert.Nil(t, err)
}
