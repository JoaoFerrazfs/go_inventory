package infrastructure

import (
	"strings"

	errors "go_inventory/Helpers/Errors"
	domain "go_inventory/SupplyInventory/Domain"

	"gorm.io/gorm"
)

type PalletRackRepository interface {
	Create(name string, location string, totalCapacity int) (*domain.PalletRackEntity, error)
	ListRacks() ([]domain.PalletRackEntity, error)
	FindPalletById(id uint) (*domain.PalletRackEntity, *errors.AppError)
	DeleteRack(id uint) (bool, *errors.AppError)
}

type palletRackRepository struct {
	db *gorm.DB
}

func NewPalletRackRepository(db *gorm.DB) PalletRackRepository {
	return &palletRackRepository{db: db}
}

func (repository *palletRackRepository) Create(name string, location string, totalCapacity int) (*domain.PalletRackEntity, error) {
	palletRack := domain.NewPalletRackEntity(name, location, totalCapacity)

	if err := repository.db.Save(&palletRack).Error; err != nil {
		return nil, err
	}

	return palletRack, nil
}

func (repository *palletRackRepository) ListRacks() ([]domain.PalletRackEntity, error) {
	var racks []domain.PalletRackEntity

	if err := repository.db.Preload("Pallets").Find(&racks).Error; err != nil {
		return nil, err
	}

	return racks, nil
}

func (repository *palletRackRepository) FindPalletById(id uint) (*domain.PalletRackEntity, *errors.AppError) {
	var rack domain.PalletRackEntity

	if err := repository.db.Preload("Pallets.PalletizedProduct").
		Preload("Pallets").First(&rack, id).Error; err != nil {
		return nil, errors.NewAppError("Rack not found", 404)
	}

	return &rack, nil
}

func (repository *palletRackRepository) DeleteRack(id uint) (bool, *errors.AppError) {
	var rack domain.PalletRackEntity

	result := repository.db.Delete(&rack, id)

	if result.Error != nil {

		if strings.Contains(result.Error.Error(), "Error 1451") {
			return false, errors.NewAppError(
				"Cannot delete rack: it has pallets associated",
				422,
			)
		}

		return false, errors.NewAppError(result.Error.Error(), 500)
	}

	if result.RowsAffected == 0 {
		return false, errors.NewAppError("Rack not found", 404)
	}

	return true, nil
}
