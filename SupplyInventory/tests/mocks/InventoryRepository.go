package mocks

import (
	appErrors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"

	"github.com/stretchr/testify/mock"
)

type InventoryRepository struct {
	mock.Mock
}

func (m *InventoryRepository) Exists(id uint) (bool, *appErrors.AppError) {
	args := m.Called(id)
	return args.Bool(0), args.Get(1).(*appErrors.AppError)
}

func (m *InventoryRepository) Where(conditions ...map[string]any) ([]entities.InventoryEntity, *appErrors.AppError) {
	args := m.Called(conditions)
	result := args.Get(0)
	if result == nil {
		return nil, args.Get(1).(*appErrors.AppError)
	}
	return result.([]entities.InventoryEntity), args.Get(1).(*appErrors.AppError)
}

func (m *InventoryRepository) FindById(id uint) (*entities.InventoryEntity, *appErrors.AppError) {
	args := m.Called(id)
	result := args.Get(0)
	if result == nil {
		return nil, args.Get(1).(*appErrors.AppError)
	}
	return result.(*entities.InventoryEntity), args.Get(1).(*appErrors.AppError)
}

func (m *InventoryRepository) Create(inventory *entities.InventoryEntity) *appErrors.AppError {
	args := m.Called(inventory)
	return args.Get(0).(*appErrors.AppError)
}

func (m *InventoryRepository) Update(inventory *entities.InventoryEntity) *appErrors.AppError {
	args := m.Called(inventory)
	return args.Get(0).(*appErrors.AppError)
}
