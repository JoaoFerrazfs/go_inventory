package infrastructure

import (
	errors "go_inventory/Helpers/Errors"
	domain "go_inventory/SupplyInventory/Domain/Entities"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	repositories "go_inventory/SupplyInventory/Domain/contracts/repositories/PalletRack"
	"strings"

	dbadapter "go_inventory/SupplyInventory/Infrastructure/repositories/db"
)

type palletRackRepository struct {
	db dbadapter.DBAdapter
}

func NewPalletRackRepository(db dbadapter.DBAdapter) repositories.PalletRackRepository {
	return &palletRackRepository{db: db}
}

func (repository *palletRackRepository) Create(name string, location string, totalCapacity int, inventoryID uint) (*entities.PalletRackEntity, error) {
	palletRack := entities.NewPalletRackEntity(inventoryID, name, location, totalCapacity)

	if err := repository.db.Save(&palletRack); err != nil {
		return nil, err
	}

	return palletRack, nil
}

func (repository *palletRackRepository) ListRacks(inventoryID *uint, page int, limit int) ([]entities.PalletRackEntity, int64, error) {
	var racks []entities.PalletRackEntity
	var total int64

	where := ""
	var args []interface{}
	if inventoryID != nil {
		where = "inventory_id = ?"
		args = append(args, *inventoryID)
	}

	err := repository.db.CountAndPaginatedFind(&entities.PalletRackEntity{}, &racks, &total, page, limit, []string{"Pallets"}, where, args...)
	if err != nil {
		return nil, 0, err
	}

	return racks, total, nil
}

func (repository *palletRackRepository) FindPalletById(id uint) (*entities.PalletRackEntity, *errors.AppError) {
	var rack entities.PalletRackEntity

	if err := repository.db.PreloadFind(&rack, "Pallets.PalletizedProduct", id); err != nil {
		return nil, errors.NewAppError("Rack not found", 404)
	}

	return &rack, nil
}

func (repository *palletRackRepository) DeleteRack(id uint) (bool, *errors.AppError) {
	var rack domain.PalletRackEntity

	rows, err := repository.db.DeleteByID(&rack, id)

	if err != nil {

		if strings.Contains(err.Error(), "Error 1451") {
			return false, errors.NewAppError(
				"Cannot delete rack: it has pallets associated",
				422,
			)
		}

		return false, errors.NewAppError(err.Error(), 500)
	}

	if rows == 0 {
		return false, errors.NewAppError("Rack not found", 404)
	}

	return true, nil
}
