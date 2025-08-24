package infrastructure

import (
	domain "go_inventory/SupplyInventory/Domain"

	"gorm.io/gorm"
)

type PalletRackRepository interface {
	Create(name string) (*domain.PalletRackEntity, error)
	ListRacks() ([]domain.PalletRackEntity, error)
}

type palletRackRepository struct {
	db *gorm.DB
}

func NewPalletRackRepository(db *gorm.DB) PalletRackRepository {
	return &palletRackRepository{db: db}
}

func (repository *palletRackRepository) Create(name string) (*domain.PalletRackEntity, error) {
	var palletRack domain.PalletRackEntity
	palletRack.Name = name
	if err := repository.db.Save(&palletRack).Error; err != nil {
		return nil, err
	}

	return &palletRack, nil
}

func (repository *palletRackRepository) ListRacks() ([]domain.PalletRackEntity, error) {
	var racks []domain.PalletRackEntity

	if err := repository.db.Preload("Pallets").Find(&racks).Error; err != nil {
		return nil, err
	}

	return racks, nil
}
