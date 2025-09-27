package services

import (
	errors "go_inventory/Helpers/Errors"
	apiContracts "go_inventory/SupplyInventory/Application/ApiContracts"
	domain "go_inventory/SupplyInventory/Domain"
	infrastructure "go_inventory/SupplyInventory/Infrastructure"
)

type PalletRackService interface {
	Create(name string, location string, totalCapacity int) (*domain.PalletRackEntity, error)
	ListRacks() ([]apiContracts.TransformedRack, error)
	FindPalletById(id uint) (*domain.PalletRackEntity, *errors.AppError)
	DeleteRack(id uint) (bool, *errors.AppError)
}

type palletRackService struct {
	repository infrastructure.PalletRackRepository
}

func NewPalletRackService(repository infrastructure.PalletRackRepository) PalletRackService {
	return &palletRackService{repository: repository}
}

func (service *palletRackService) Create(name string, location string, totalCapacity int) (*domain.PalletRackEntity, error) {
	newPalletRack, err := service.repository.Create(name, location, totalCapacity)
	if err != nil {
		return nil, err
	}

	return newPalletRack, nil
}

func (service *palletRackService) ListRacks() ([]apiContracts.TransformedRack, error) {
	racks, err := service.repository.ListRacks()
	if err != nil {
		return nil, err
	}

	newIndices := []apiContracts.TransformedRack{}

	for _, valor := range racks {

		transformedRack := apiContracts.TransformedRack{
			ID:            valor.ID,
			Name:          valor.Name,
			Pallets:       valor.Pallets,
			Location:      valor.Location,
			TotalCapacity: valor.TotalCapacity,
			PercetageUsed: (float64(len(valor.Pallets)) / float64(valor.TotalCapacity)) * 100,
		}

		newIndices = append(newIndices, transformedRack)
	}

	return newIndices, nil
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
