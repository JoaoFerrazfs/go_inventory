package services

import (
	domain "go_inventory/SupplyInventory/Domain"
	infrastructure "go_inventory/SupplyInventory/Infrastructure"
)

type PalletRackService interface {
	Create(palletRack domain.PalletRackEntity) (*domain.PalletRackEntity, error)
}

type palletRackService struct {
	repository infrastructure.PalletRackRepository
}

func NewPalletRackService(repository infrastructure.PalletRackRepository) PalletRackService {
	return &palletRackService{repository: repository}
}

func (service *palletRackService) Create(palletRack domain.PalletRackEntity) (*domain.PalletRackEntity, error) {
	newPalletRack, err := service.repository.Create(palletRack)
	if err != nil {
		return nil, err
	}

	return newPalletRack, nil
}
