package inventory

import (
	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	repositories "go_inventory/SupplyInventory/Domain/contracts/repositories/Inventory"
	"time"
)

type InventoryService interface {
	ListInventories() ([]entities.InventoryEntity, *errors.AppError)
	GetInventoryByID(id uint) (entities.InventoryEntity, *errors.AppError)
	CreateInventory(name string, description string, user entities.UserEntity) (entities.InventoryEntity, *errors.AppError)
	UpdateInventory(id uint, data interface{}) (entities.InventoryEntity, *errors.AppError)
}

type inventoryService struct {
	inventoryRepository repositories.InventoryRepository
}

func (service *inventoryService) ListInventories() ([]entities.InventoryEntity, *errors.AppError) {
	inventories, appError := service.inventoryRepository.Where()
	if appError != nil {
		return nil, appError
	}
	return inventories, nil
}

func (service *inventoryService) GetInventoryByID(id uint) (entities.InventoryEntity, *errors.AppError) {
	inventory, appErr := service.inventoryRepository.FindById(id)
	if appErr != nil {
		return entities.InventoryEntity{}, appErr
	}
	return *inventory, nil
}

func (service *inventoryService) CreateInventory(name string, description string, user entities.UserEntity) (entities.InventoryEntity, *errors.AppError) {
	inventory := entities.InventoryEntity{}
	inventory.User = user
	inventory.UserID = user.ID
	inventory.StartedAt = time.Now()
	inventory.Name = name
	inventory.Description = description

	err := service.inventoryRepository.Create(&inventory)
	if err != nil {
		return entities.InventoryEntity{}, err
	}

	return inventory, nil
}

func (service *inventoryService) UpdateInventory(id uint, data interface{}) (entities.InventoryEntity, *errors.AppError) {
	// TODO: Implement logic to update inventory
	return entities.InventoryEntity{}, nil
}

func NewInventoryService(
	repository repositories.InventoryRepository,
) InventoryService {
	return &inventoryService{
		inventoryRepository: repository,
	}
}
