package inventory

import (
	"testing"

	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockInventoryRepo struct{ mock.Mock }

func (m *mockInventoryRepo) Exists(id uint) (bool, *errors.AppError) {
	args := m.Called(id)
	return args.Bool(0), args.Get(1).(*errors.AppError)
}

func (m *mockInventoryRepo) List() ([]entities.InventoryEntity, *errors.AppError) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).([]entities.InventoryEntity), nil
}

func (m *mockInventoryRepo) Where(conditions map[string]interface{}) ([]entities.InventoryEntity, *errors.AppError) {
	args := m.Called(conditions)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).([]entities.InventoryEntity), nil
}

func (m *mockInventoryRepo) FindById(id uint) (*entities.InventoryEntity, *errors.AppError) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*entities.InventoryEntity), nil
}

func (m *mockInventoryRepo) Create(inventory *entities.InventoryEntity) *errors.AppError {
	args := m.Called(inventory)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*errors.AppError)
}

func (m *mockInventoryRepo) Update(inventory *entities.InventoryEntity) *errors.AppError {
	args := m.Called(inventory)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*errors.AppError)
}

func TestListInventories_Success(t *testing.T) {
	// Set
	repo := &mockInventoryRepo{}
	svc := NewInventoryService(repo)
	expectedInventories := []entities.InventoryEntity{
		{ID: 1, UserID: 1, Status: entities.InventoryStatusOpen},
		{ID: 2, UserID: 2, Status: entities.InventoryStatusClosed},
	}

	// Expectations
	repo.On("List").Return(expectedInventories, nil)

	// Actions
	inventories, err := svc.ListInventories()

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, inventories)
	assert.Equal(t, 2, len(inventories))
	assert.Equal(t, expectedInventories, inventories)
	repo.AssertExpectations(t)
}

func TestListInventories_EmptyList(t *testing.T) {
	// Set
	repo := &mockInventoryRepo{}
	svc := NewInventoryService(repo)
	expectedInventories := []entities.InventoryEntity{}

	// Expectations
	repo.On("List").Return(expectedInventories, nil)

	// Actions
	inventories, err := svc.ListInventories()

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, inventories)
	assert.Equal(t, 0, len(inventories))
	repo.AssertExpectations(t)
}

func TestListInventories_Error(t *testing.T) {
	// Set
	repo := &mockInventoryRepo{}
	svc := NewInventoryService(repo)
	appErr := errors.NewAppError("Failed to list inventories", 500)

	// Expectations
	repo.On("List").Return(nil, appErr)

	// Actions
	inventories, err := svc.ListInventories()

	// Assertions
	assert.Nil(t, inventories)
	assert.NotNil(t, err)
	assert.Equal(t, appErr, err)
	repo.AssertExpectations(t)
}
