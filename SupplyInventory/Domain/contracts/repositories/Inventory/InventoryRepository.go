package Inventory

import (
	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
)

type InventoryRepository interface {
	Exists(id uint) (bool, *errors.AppError)
	Where(conditions map[string]interface{}) ([]entities.InventoryEntity, *errors.AppError)
	FindById(id uint) (*entities.InventoryEntity, *errors.AppError)
	Create(inventory *entities.InventoryEntity) *errors.AppError
	Update(inventory *entities.InventoryEntity) *errors.AppError
}
