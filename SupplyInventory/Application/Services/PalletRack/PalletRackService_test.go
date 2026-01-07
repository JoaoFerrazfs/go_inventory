package services

import (
	"testing"

	errors "go_inventory/Helpers/Errors"
	apiContracts "go_inventory/SupplyInventory/Application/ApiContracts"
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

func (m *mockPalletRackRepo) ListRacks() ([]entities.PalletRackEntity, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.PalletRackEntity), args.Error(1)
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
	svc := NewPalletRackService(repo)

	// Expectations
	created := &entities.PalletRackEntity{ID: 1, Name: "R1", Location: "L1", TotalCapacity: 10}
	repo.On("Create", "R1", "L1", 10, uint(0)).Return(created, nil)

	// Actions
	res, err := svc.Create("R1", "L1", 10, uint(0))

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, created, res)
	repo.AssertExpectations(t)
}

func TestListRacks_TransformsPercentage(t *testing.T) {
	// Set
	repo := &mockPalletRackRepo{}
	svc := NewPalletRackService(repo)

	// Expectations
	racks := []entities.PalletRackEntity{
		{ID: 1, Name: "R1", Location: "L1", TotalCapacity: 4, Pallets: []entities.PalletEntity{{}, {}}},
	}
	repo.On("ListRacks").Return(racks, nil)

	// Actions
	res, err := svc.ListRacks()

	// Assertions
	assert.Nil(t, err)
	if len(res) > 0 {
		expected := apiContracts.TransformedRack{ID: 1, Name: "R1", Location: "L1", TotalCapacity: 4, PercetageUsed: 50.0, Pallets: racks[0].Pallets}
		assert.Equal(t, expected, res[0])
	}
	repo.AssertExpectations(t)
}

func TestFindPalletById_RepoError(t *testing.T) {
	// Set
	repo := &mockPalletRackRepo{}
	svc := NewPalletRackService(repo)

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
	svc := NewPalletRackService(repo)

	// Expectations
	repo.On("DeleteRack", uint(5)).Return(true, nil)

	// Actions
	ok, err := svc.DeleteRack(5)

	// Assertions
	assert.Nil(t, err)
	assert.True(t, ok)
	repo.AssertExpectations(t)
}
