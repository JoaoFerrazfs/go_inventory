package services

import (
	domain "go_inventory/SupplyInventory/Domain"
	infrastructure "go_inventory/SupplyInventory/Infrastructure"
)

type PalletizedProductService interface {
	AddProductsToPallet(PalletID uint, Ean int, Quantity int) (*domain.PalletEntity, error)
	DeleteProductsFromPallet(palletId uint, productsEan int) (bool, error)
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

func (service *palletizedProductService) AddProductsToPallet(PalletID uint, Ean int, Quantity int) (*domain.PalletEntity, error) {
	product := domain.PalletizedProductEntity{
		PalletID: PalletID,
		EAN:      Ean,
		Quantity: Quantity,
	}

	_, err := service.palletizedProductRepository.AddProductsToPallet(product)
	if err != nil {
		return nil, err
	}

	pallet, err := service.palletRepository.GetSupplyById(product.PalletID)
	if err != nil {
		return nil, err
	}

	return pallet, nil
}

func (service *palletizedProductService) DeleteProductsFromPallet(palletId uint, productsEan int) (bool, error) {
	deleted, err := service.palletizedProductRepository.DeleteProductsFromPallet(palletId, productsEan)
	if err != nil {
		return false, err
	}

	return deleted, nil
}
