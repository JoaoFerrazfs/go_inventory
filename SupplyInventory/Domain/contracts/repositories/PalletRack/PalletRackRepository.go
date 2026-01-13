package repositories

import (
	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
)

type PalletRackRepository interface {
	Create(name string, location string, totalCapacity int, inventoryID uint) (*entities.PalletRackEntity, error)
	ListRacks(inventoryID *uint, page int, limit int) ([]entities.PalletRackEntity, int64, error)
	FindPalletById(id uint) (*entities.PalletRackEntity, *errors.AppError)
	DeleteRack(id uint) (bool, *errors.AppError)
}
