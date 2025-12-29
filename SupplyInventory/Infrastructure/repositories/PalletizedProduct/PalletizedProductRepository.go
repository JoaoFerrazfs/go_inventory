package infrastructure

import (
	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	palletRepo "go_inventory/SupplyInventory/Domain/contracts/repositories/Pallet"
	productRepo "go_inventory/SupplyInventory/Domain/contracts/repositories/PalletizedProduct"
	dbadapter "go_inventory/SupplyInventory/Infrastructure/repositories/db"
)

type PalletizedProductRepositoryImpl struct {
	db               dbadapter.DBAdapter
	palletRepository palletRepo.PalletRepository
}

func NewPalletizedProductRepository(db dbadapter.DBAdapter, palletRepository palletRepo.PalletRepository) productRepo.PalletizedProductRepository {
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

			if err := repository.db.Save(&pallet.PalletizedProduct[i]); err != nil {
				return false, errors.NewAppError(err.Error(), 500)
			}

			return true, nil
		}
	}

	if err := repository.db.AppendAssociation(pallet, &product); err != nil {
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
			rows, err := repository.db.DeleteByID(&pallet.PalletizedProduct[i], pallet.PalletizedProduct[i].ID)
			if err != nil {
				return false, errors.NewAppError(err.Error(), 500)
			}
			if rows == 0 {
				return false, errors.NewAppError("product not deleted", 500)
			}
			return true, nil
		}
	}

	return false, errors.NewAppError("product is not in the pallets", 404)
}
