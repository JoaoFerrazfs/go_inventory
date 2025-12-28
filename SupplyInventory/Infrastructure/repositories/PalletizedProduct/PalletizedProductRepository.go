package infrastructure

import (
	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	palletRepo "go_inventory/SupplyInventory/Domain/contracts/repositories/Pallet"
	productRepo "go_inventory/SupplyInventory/Domain/contracts/repositories/PalletizedProduct"

	"gorm.io/gorm"
)

type PalletizedProductRepositoryImpl struct {
       db               *gorm.DB
       palletRepository palletRepo.PalletRepository
}

func NewPalletizedProductRepository(db *gorm.DB, palletRepository palletRepo.PalletRepository) productRepo.PalletizedProductRepository {
	return &PalletizedProductRepositoryImpl{db: db, palletRepository: palletRepository}
}

func (repository *PalletizedProductRepositoryImpl) AddProductsToPallet(product entities.PalletizedProductEntity) (bool, *errors.AppError) {
	pallet, err := repository.palletRepository.GetSupplyById(product.PalletID)
	if err != nil || pallet == nil {
		return false, err
	}

	for i := range pallet.PalletizedProduct {
		if pallet.PalletizedProduct[i].EAN == product.EAN {
			pallet.PalletizedProduct[i].Quantity = product.Quantity

			if err := repository.db.Save(&pallet.PalletizedProduct[i]).Error; err != nil {
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