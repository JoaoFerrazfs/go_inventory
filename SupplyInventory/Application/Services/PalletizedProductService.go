package services

import (
	domain "go_inventory/SupplyInventory/Domain"
	infrastructure "go_inventory/SupplyInventory/Infrastructure"
)

type PalletizedProductService interface {
	AddProductsToPallet(product domain.PalletizedProductEntity) (*domain.PalletEntity, error)
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

func (service *palletizedProductService) AddProductsToPallet(product domain.PalletizedProductEntity) (*domain.PalletEntity, error) {
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
