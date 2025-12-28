package repositories

import (
	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
)

type PalletRackRepository interface {
	Create(name string, location string, totalCapacity int) (*entities.PalletRackEntity, error)
	ListRacks() ([]entities.PalletRackEntity, error)
	FindPalletById(id uint) (*entities.PalletRackEntity, *errors.AppError)
	DeleteRack(id uint) (bool, *errors.AppError)
}
