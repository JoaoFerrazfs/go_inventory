package services

import (
	errors "go_inventory/Helpers/Errors"
	apiContracts "go_inventory/SupplyInventory/Application/ApiContracts"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	repositoriesPalletRack "go_inventory/SupplyInventory/Domain/contracts/repositories/PalletRack"
)

type PalletRackService interface {
	Create(name string, location string, totalCapacity int, inventoryID uint) (*entities.PalletRackEntity, error)
	ListRacks(inventoryID *uint, page int, limit int) (*apiContracts.PaginatedRacksResponse, error)
	FindPalletById(id uint) (*entities.PalletRackEntity, *errors.AppError)
	DeleteRack(id uint) (bool, *errors.AppError)
}

type palletRackService struct {
	repository repositoriesPalletRack.PalletRackRepository
}

func NewPalletRackService(repository repositoriesPalletRack.PalletRackRepository) PalletRackService {
	return &palletRackService{repository: repository}
}

func (service *palletRackService) Create(name string, location string, totalCapacity int, inventoryID uint) (*entities.PalletRackEntity, error) {
	newPalletRack, err := service.repository.Create(name, location, totalCapacity, inventoryID)
	if err != nil {
		return nil, err
	}

	return newPalletRack, nil
}

func (service *palletRackService) ListRacks(inventoryID *uint, page int, limit int) (*apiContracts.PaginatedRacksResponse, error) {
	racks, total, err := service.repository.ListRacks(inventoryID, page, limit)
	if err != nil {
		return nil, err
	}

	transformedRacks := []apiContracts.TransformedRack{}

	for _, rack := range racks {
		transformedRack := apiContracts.TransformedRack{
			ID:            rack.ID,
			Name:          rack.Name,
			Pallets:       rack.Pallets,
			Location:      rack.Location,
			TotalCapacity: rack.TotalCapacity,
			PercetageUsed: (float64(len(rack.Pallets)) / float64(rack.TotalCapacity)) * 100,
		}
		transformedRacks = append(transformedRacks, transformedRack)
	}

	response := &apiContracts.PaginatedRacksResponse{
		Data:  transformedRacks,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return response, nil
}

func (service *palletRackService) FindPalletById(id uint) (*entities.PalletRackEntity, *errors.AppError) {
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
