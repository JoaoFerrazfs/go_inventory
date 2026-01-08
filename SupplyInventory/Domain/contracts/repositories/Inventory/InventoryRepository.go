package Inventory

import errors "go_inventory/Helpers/Errors"

type InventoryRepository interface {
	Exists(id uint) (bool, *errors.AppError)
}
