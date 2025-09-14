package services

import (
	errors "go_inventory/Helpers/Errors"
	domain "go_inventory/SupplyInventory/Domain"
	infrastructure "go_inventory/SupplyInventory/Infrastructure"
)

type PalletRackService interface {
	Create(name string) (*domain.PalletRackEntity, error)
	ListRacks() ([]domain.PalletRackEntity, error)
	FindPalletById(id uint) (*domain.PalletRackEntity, *errors.AppError)
	DeleteRack(id uint) (bool, *errors.AppError)
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

func (service *palletRackService) FindPalletById(id uint) (*domain.PalletRackEntity, *errors.AppError) {
	racks, appErr := service.repository.FindPalletById(id)
	if appErr != nil {
		return nil, appErr
	}

	return racks, nil
}

func (service *palletRackService) DeleteRack(id uint) (bool, *errors.AppError) {
	_, appErr := service.repository.DeleteRack(id)
	if appErr != nil {
		return false, appErr
	}

	return true, nil
}
