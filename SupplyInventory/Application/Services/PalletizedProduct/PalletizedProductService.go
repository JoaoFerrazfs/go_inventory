// (removed invalid import lines)
package services

import (
	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	repositories "go_inventory/SupplyInventory/Domain/contracts/repositories/Pallet"
	repositoriesProduct "go_inventory/SupplyInventory/Domain/contracts/repositories/PalletizedProduct"
)

type PalletizedProductService interface {
	AddProductsToPallet(PalletID uint, Ean int, Quantity int, InventoryID uint) (*entities.PalletEntity, *errors.AppError)
	DeleteProductsFromPallet(palletId uint, productsEan int) (bool, *errors.AppError)
}

type palletizedProductService struct {
	palletRepository            repositories.PalletRepository
	palletizedProductRepository repositoriesProduct.PalletizedProductRepository
}

func NewPalletizedProductService(
	palletRepository repositories.PalletRepository,
	palletizedProductRepository repositoriesProduct.PalletizedProductRepository,
) PalletizedProductService {
	return &palletizedProductService{palletRepository: palletRepository, palletizedProductRepository: palletizedProductRepository}
}

func (service *palletizedProductService) AddProductsToPallet(PalletID uint, Ean int, Quantity int, InventoryID uint) (*entities.PalletEntity, *errors.AppError) {
	pallet, appErr := service.palletRepository.GetSupplyById(PalletID)
	if appErr != nil {
		return nil, appErr
	}

	if pallet.InventoryID != InventoryID {
		return nil, errors.NewAppError("Pallet does not belong to the same inventory", 422)
	}

	product := entities.PalletizedProductEntity{
		PalletID:    PalletID,
		EAN:         Ean,
		Quantity:    Quantity,
		InventoryID: InventoryID,
	}

	_, err := service.palletizedProductRepository.AddProductsToPallet(product)
	if err != nil {
		return nil, errors.NewAppError(err.Error(), 400)
	}

	// Fetch updated pallet
	updatedPallet, appErr := service.palletRepository.GetSupplyById(PalletID)
	if appErr != nil {
		return nil, appErr
	}

	return updatedPallet, nil
}

func (service *palletizedProductService) DeleteProductsFromPallet(palletId uint, productsEan int) (bool, *errors.AppError) {
	deleted, appErr := service.palletizedProductRepository.DeleteProductsFromPallet(palletId, productsEan)
	if appErr != nil {
		return false, appErr
	}

	return deleted, nil
}
