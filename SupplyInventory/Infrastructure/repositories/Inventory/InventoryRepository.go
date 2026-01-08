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

func (r *inventoryRepository) Create(inv *entities.InventoryEntity) *errors.AppError {
	if err := r.db.Create(inv); err != nil {
		return errors.NewAppError("Failed to create inventory", 500)
	}
	return nil
}

func (r *inventoryRepository) FindById(id uint) (*entities.InventoryEntity, *errors.AppError) {
	var inv entities.InventoryEntity
	if err := r.db.FirstByID(&inv, id); err != nil {
		return nil, errors.NewAppError("Inventory not found", 404)
	}
	return &inv, nil
}

func (r *inventoryRepository) Where(conditions map[string]interface{}) ([]entities.InventoryEntity, *errors.AppError) {
	var inventories []entities.InventoryEntity
	if err := r.db.GetDB().Where(conditions).Find(&inventories).Error; err != nil {
		return nil, errors.NewAppError("Failed to find inventories", 404)
	}
	return inventories, nil
}

func (r *inventoryRepository) Update(inv *entities.InventoryEntity) *errors.AppError {
	if err := r.db.Save(inv); err != nil {
		return errors.NewAppError("Failed to update inventory", 500)
	}
	return nil
}
