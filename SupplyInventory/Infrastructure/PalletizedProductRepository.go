package infrastructure

import (
	domain "go_inventory/SupplyInventory/Domain"

	"gorm.io/gorm"
)

type PalletizedProductRepository interface {
	AddProductsToPallet(product domain.PalletizedProductEntity) (bool, error)
}

type PalletizedProductRepositoryImpl struct {
	db               *gorm.DB
	palletRepository PalletRepository
}

func NewPalletizedProductRepository(db *gorm.DB, palletRepository PalletRepository) PalletizedProductRepository {
	return &PalletizedProductRepositoryImpl{db: db, palletRepository: palletRepository}
}

func (repository *PalletizedProductRepositoryImpl) AddProductsToPallet(product domain.PalletizedProductEntity) (bool, error) {
	pallet, err := repository.palletRepository.GetSupplyById(product.PalletID)
	if err != nil || pallet == nil {
		return false, err
	}

	product.PalletID = pallet.ID
	if err := repository.db.Model(pallet).Association("PalletizedProduct").Append(&product); err != nil {
		return false, err
	}

	return true, nil
}
