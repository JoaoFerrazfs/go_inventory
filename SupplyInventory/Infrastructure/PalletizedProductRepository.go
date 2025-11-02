package infrastructure

import (
	errors "go_inventory/Helpers/Errors"
	domain "go_inventory/SupplyInventory/Domain"

	"gorm.io/gorm"
)

type PalletizedProductRepository interface {
	AddProductsToPallet(product domain.PalletizedProductEntity) (bool, *errors.AppError)
	DeleteProductsFromPallet(palletId uint, productsEan int) (bool, *errors.AppError)
}

type PalletizedProductRepositoryImpl struct {
	db               *gorm.DB
	palletRepository PalletRepository
}

func NewPalletizedProductRepository(db *gorm.DB, palletRepository PalletRepository) PalletizedProductRepository {
	return &PalletizedProductRepositoryImpl{db: db, palletRepository: palletRepository}
}

func (repository *PalletizedProductRepositoryImpl) AddProductsToPallet(product domain.PalletizedProductEntity) (bool, *errors.AppError) {
	pallet, err := repository.palletRepository.GetSupplyById(product.PalletID)
	if err != nil || pallet == nil {
		return false, err
	}

	for palletizedProduct := range pallet.PalletizedProduct {
		if pallet.PalletizedProduct[palletizedProduct].EAN == product.EAN {
			pallet.PalletizedProduct[palletizedProduct].Quantity = product.Quantity

			if err := repository.db.Save(&pallet.PalletizedProduct[palletizedProduct]).Error; err != nil {
				return false, errors.NewAppError(err.Error(), 500)
			}

			return true, nil
		}
	}

	if err := repository.db.Model(pallet).Association("PalletizedProduct").Append(&product); err != nil {
		return false, errors.NewAppError(err.Error(), 500)
	}

	return true, nil
}

func (repository *PalletizedProductRepositoryImpl) DeleteProductsFromPallet(palletId uint, productsEan int) (bool, *errors.AppError) {
	pallet, err := repository.palletRepository.GetSupplyById(palletId)
	if err != nil || pallet == nil {
		return false, err
	}

	for i := range pallet.PalletizedProduct {
		if pallet.PalletizedProduct[i].EAN == productsEan {
			if err := repository.db.Delete(&pallet.PalletizedProduct[i]).Error; err != nil {
				return false, errors.NewAppError(err.Error(), 500)
			}
			return true, nil
		}
	}

	return false, errors.NewAppError("product is not in the pallet", 404)
}
