package infrastructure

import (
	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	repositories "go_inventory/SupplyInventory/Domain/contracts/repositories/Inventory"

	dbadapter "go_inventory/SupplyInventory/Infrastructure/repositories/db"
)

type inventoryRepository struct {
    db dbadapter.DBAdapter
}

func NewInventoryRepository(db dbadapter.DBAdapter) repositories.InventoryRepository {
    return &inventoryRepository{db: db}
}

func (r *inventoryRepository) Exists(id uint) (bool, *errors.AppError) {
    var inv entities.InventoryEntity
    if err := r.db.FirstByID(&inv, id); err != nil {
        return false, errors.NewAppError("Inventory not found", 404)
    }
    return true, nil
}
