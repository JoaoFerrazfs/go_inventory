package inventory_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	errors "go_inventory/Helpers/Errors"
	services "go_inventory/SupplyInventory/Application/Services/Inventory"
	entities "go_inventory/SupplyInventory/Domain/Entities"
)

// Mock do InventoryRepository
type MockInventoryRepository struct {
	mock.Mock
}

func (m *MockInventoryRepository) Exists(id uint) (bool, *errors.AppError) {
	args := m.Called(id)
	return args.Bool(0), args.Get(1).(*errors.AppError)
}

func (m *MockInventoryRepository) Where(conditions ...map[string]any) ([]entities.InventoryEntity, *errors.AppError) {
	args := m.Called(conditions)
	return args.Get(0).([]entities.InventoryEntity), args.Get(1).(*errors.AppError)
}

func (m *MockInventoryRepository) FindById(id uint) (*entities.InventoryEntity, *errors.AppError) {
	args := m.Called(id)
	return args.Get(0).(*entities.InventoryEntity), args.Get(1).(*errors.AppError)
}

func (m *MockInventoryRepository) Create(inventory *entities.InventoryEntity) *errors.AppError {
	args := m.Called(inventory)
	return args.Get(0).(*errors.AppError)
}

func (m *MockInventoryRepository) Update(inventory *entities.InventoryEntity) *errors.AppError {
	args := m.Called(inventory)
	return args.Get(0).(*errors.AppError)
}

func TestCreateInventory_Success(t *testing.T) {
	// Set
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)
	user := entities.UserEntity{ID: 1, Name: "Test User"}

	mockRepo.On("Create", mock.MatchedBy(func(inv *entities.InventoryEntity) bool {
		return inv.Name == "Test Inventory" && inv.UserID == 1
	})).Return((*errors.AppError)(nil))

	// Actions
	result, err := service.CreateInventory("Test Inventory", "Test Description", user)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, "Test Inventory", result.Name)
	assert.Equal(t, uint(1), result.UserID)
	mockRepo.AssertExpectations(t)
}

func TestListInventories_Success(t *testing.T) {
	// Set
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)
	inventories := []entities.InventoryEntity{
		{ID: 1, Name: "Inventory 1", User: entities.UserEntity{ID: 1, Name: "User 1"}},
		{ID: 2, Name: "Inventory 2", User: entities.UserEntity{ID: 2, Name: "User 2"}},
	}

	mockRepo.On("Where", mock.Anything).Return(inventories, (*errors.AppError)(nil))

	// Actions
	result, err := service.ListInventories()

	// Assertions
	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Inventory 1", result[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestGetInventoryByID_Success(t *testing.T) {
	// Set
	mockRepo := new(MockInventoryRepository)
	service := services.NewInventoryService(mockRepo)
	inventory := &entities.InventoryEntity{
		ID:   1,
		Name: "Test Inventory",
		User: entities.UserEntity{ID: 1, Name: "Test User"},
	}

	mockRepo.On("FindById", uint(1)).Return(inventory, (*errors.AppError)(nil))

	// Actions
	result, err := service.GetInventoryByID(1)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Test Inventory", result.Name)
	mockRepo.AssertExpectations(t)
}
