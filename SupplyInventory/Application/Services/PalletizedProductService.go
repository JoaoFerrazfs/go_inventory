package services

import (
	errors "go_inventory/Helpers/Errors"
	domain "go_inventory/SupplyInventory/Domain"
	infrastructure "go_inventory/SupplyInventory/Infrastructure"
)

type PalletizedProductService interface {
	AddProductsToPallet(PalletID uint, Ean int, Quantity int) (*domain.PalletEntity, *errors.AppError)
	DeleteProductsFromPallet(palletId uint, productsEan int) (bool, *errors.AppError)
}

type palletizedProductService struct {
	palletRepository            infrastructure.PalletRepository
	palletizedProductRepository infrastructure.PalletizedProductRepository
}

func NewPalletizedProductService(
	palletRepository infrastructure.PalletRepository,
	palletizedProductRepository infrastructure.PalletizedProductRepository,
) PalletizedProductService {
	return &palletizedProductService{palletRepository: palletRepository, palletizedProductRepository: palletizedProductRepository}
}

func (service *palletizedProductService) AddProductsToPallet(PalletID uint, Ean int, Quantity int) (*domain.PalletEntity, *errors.AppError) {
	product := domain.PalletizedProductEntity{
		PalletID: PalletID,
		EAN:      Ean,
		Quantity: Quantity,
	}

	_, err := service.palletizedProductRepository.AddProductsToPallet(product)
	if err != nil {
		return nil, errors.NewAppError(err.Error(), 400)
	}

	pallet, appErr := service.palletRepository.GetSupplyById(product.PalletID)
	if appErr != nil {
		return nil, appErr
	}

	return pallet, nil
}

func (service *palletizedProductService) DeleteProductsFromPallet(palletId uint, productsEan int) (bool, *errors.AppError) {
	deleted, appErr := service.palletizedProductRepository.DeleteProductsFromPallet(palletId, productsEan)
	if appErr != nil {
		return false, appErr
	}

	return deleted, nil
}
