package services_test

import (
	"testing"

	errors "go_inventory/Helpers/Errors"
	apiContracts "go_inventory/SupplyInventory/Application/ApiContracts"
	palletRackService "go_inventory/SupplyInventory/Application/Services/PalletRack"
	entities "go_inventory/SupplyInventory/Domain/Entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPalletRackRepo struct{ mock.Mock }

func (m *mockPalletRackRepo) Create(name string, location string, totalCapacity int, inventoryID uint) (*entities.PalletRackEntity, error) {
	args := m.Called(name, location, totalCapacity, inventoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.PalletRackEntity), args.Error(1)
}

func (m *mockPalletRackRepo) ListRacks(inventoryID *uint, page int, limit int) ([]entities.PalletRackEntity, int64, error) {
	args := m.Called(inventoryID, page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]entities.PalletRackEntity), args.Get(1).(int64), args.Error(2)
}

func (m *mockPalletRackRepo) FindPalletById(id uint) (*entities.PalletRackEntity, *errors.AppError) {
	args := m.Called(id)
	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil
		}
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*entities.PalletRackEntity), nil
}

func (m *mockPalletRackRepo) DeleteRack(id uint) (bool, *errors.AppError) {
	args := m.Called(id)
	if args.Get(1) != nil {
		return args.Bool(0), args.Get(1).(*errors.AppError)
	}
	return args.Bool(0), nil
}

func TestCreate_Success(t *testing.T) {
	// Set
	repo := &mockPalletRackRepo{}
	svc := palletRackService.NewPalletRackService(repo)

	// Expectations
	created := &entities.PalletRackEntity{ID: 1, Name: "R1", Location: "L1", TotalCapacity: 10, InventoryID: 1}
	repo.On("Create", "R1", "L1", 10, uint(1)).Return(created, nil)

	// Actions
	res, err := svc.Create("R1", "L1", 10, uint(1))

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, created, res)
	repo.AssertExpectations(t)
}

func TestListRacks_TransformsPercentage(t *testing.T) {
	// Set
	repo := &mockPalletRackRepo{}
	svc := palletRackService.NewPalletRackService(repo)

	// Expectations
	racks := []entities.PalletRackEntity{
		{ID: 1, Name: "R1", Location: "L1", TotalCapacity: 4, InventoryID: 1, Pallets: []entities.PalletEntity{{}, {}}},
	}
	inventoryID := uint(1)
	repo.On("ListRacks", &inventoryID, 1, 10).Return(racks, int64(len(racks)), nil)

	// Actions
	res, err := svc.ListRacks(&inventoryID, 1, 10)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, int64(1), res.Total)
	assert.Equal(t, 1, res.Page)
	assert.Equal(t, 10, res.Limit)
	if len(res.Data) > 0 {
		expected := apiContracts.TransformedRack{ID: 1, Name: "R1", Location: "L1", TotalCapacity: 4, PercetageUsed: 50.0, Pallets: racks[0].Pallets}
		assert.Equal(t, expected, res.Data[0])
	}
	repo.AssertExpectations(t)
}

func TestListRacks_Admin_NoInventoryID(t *testing.T) {
	// Set
	repo := &mockPalletRackRepo{}
	svc := palletRackService.NewPalletRackService(repo)

	// Expectations
	racks := []entities.PalletRackEntity{
		{ID: 1, Name: "R1", Location: "L1", TotalCapacity: 4, InventoryID: 1, Pallets: []entities.PalletEntity{{}, {}}},
		{ID: 2, Name: "R2", Location: "L2", TotalCapacity: 10, InventoryID: 2, Pallets: []entities.PalletEntity{}},
	}
	repo.On("ListRacks", (*uint)(nil), 1, 20).Return(racks, int64(len(racks)), nil)

	// Actions
	res, err := svc.ListRacks(nil, 1, 20)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, int64(2), res.Total)
	assert.Equal(t, 2, len(res.Data))
	assert.Equal(t, 20, res.Limit)
	repo.AssertExpectations(t)
}

func TestFindPalletById_RepoError(t *testing.T) {
	// Set
	repo := &mockPalletRackRepo{}
	svc := palletRackService.NewPalletRackService(repo)

	// Expectations
	appErr := errors.NewAppError("not found", 404)
	repo.On("FindPalletById", uint(99)).Return(nil, appErr)

	// Actions
	res, err := svc.FindPalletById(99)

	// Assertions
	assert.Nil(t, res)
	assert.Equal(t, appErr, err)
	repo.AssertExpectations(t)
}

func TestDeleteRack_Success(t *testing.T) {
	// Set
	repo := &mockPalletRackRepo{}
	svc := palletRackService.NewPalletRackService(repo)

	// Expectations
	repo.On("DeleteRack", uint(5)).Return(true, nil)

	// Actions
	ok, err := svc.DeleteRack(5)

	// Assertions
	assert.Nil(t, err)
	assert.True(t, ok)
	repo.AssertExpectations(t)
}
