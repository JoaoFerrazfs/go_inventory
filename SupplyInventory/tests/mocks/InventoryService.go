package mocks

import (
	appErrors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"

	"github.com/stretchr/testify/mock"
)

type InventoryService struct {
	mock.Mock
}

func (m *InventoryService) ListInventories() ([]entities.InventoryEntity, *appErrors.AppError) {
	args := m.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Get(1).(*appErrors.AppError)
	}
	return result.([]entities.InventoryEntity), args.Get(1).(*appErrors.AppError)
}

func (m *InventoryService) GetInventoryByID(id uint) (entities.InventoryEntity, *appErrors.AppError) {
	args := m.Called(id)
	return args.Get(0).(entities.InventoryEntity), args.Get(1).(*appErrors.AppError)
}

func (m *InventoryService) CreateInventory(name string, description string, user entities.UserEntity) (entities.InventoryEntity, *appErrors.AppError) {
	args := m.Called(name, description, user)
	return args.Get(0).(entities.InventoryEntity), args.Get(1).(*appErrors.AppError)
}

func (m *InventoryService) UpdateInventory(id uint, data interface{}) (entities.InventoryEntity, *appErrors.AppError) {
	args := m.Called(id, data)
	return args.Get(0).(entities.InventoryEntity), args.Get(1).(*appErrors.AppError)
}