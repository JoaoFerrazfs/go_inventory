package Inventory

import errors "go_inventory/Helpers/Errors"

type InventoryRepository interface {
	Exists(id uint) (bool, *errors.AppError)
<<<<<<< Updated upstream
=======
	Where(conditions ...map[string]any) ([]entities.InventoryEntity, *errors.AppError)
	FindById(id uint) (*entities.InventoryEntity, *errors.AppError)
	Create(inventory *entities.InventoryEntity) *errors.AppError
	Update(inventory *entities.InventoryEntity) *errors.AppError
>>>>>>> Stashed changes
}
