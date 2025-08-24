package services

import (
	domain "go_inventory/SupplyInventory/Domain"
	infrastructure "go_inventory/SupplyInventory/Infrastructure"
)

type PalletRackService interface {
	Create(name string) (*domain.PalletRackEntity, error)
	ListRacks() ([]domain.PalletRackEntity, error)
}

type palletRackService struct {
	repository infrastructure.PalletRackRepository
}

func NewPalletRackService(repository infrastructure.PalletRackRepository) PalletRackService {
	return &palletRackService{repository: repository}
}

func (service *palletRackService) Create(name string) (*domain.PalletRackEntity, error) {
	newPalletRack, err := service.repository.Create(name)
	if err != nil {
		return nil, err
	}

	return newPalletRack, nil
}

func (service *palletRackService) ListRacks() ([]domain.PalletRackEntity, error) {
	racks, err := service.repository.ListRacks()
	if err != nil {
		return nil, err
	}

	return racks, nil
}
