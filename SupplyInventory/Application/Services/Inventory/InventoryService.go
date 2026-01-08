package inventory

import (
	entities "go_inventory/SupplyInventory/Domain/Entities"
	repositories "go_inventory/SupplyInventory/Domain/contracts/repositories/Inventory"
)

type InventoryService interface {
	ListInventories() ([]entities.InventoryEntity, error)
	GetInventoryByID(id uint) (entities.InventoryEntity, error)
	CreateInventory(data interface{}) (entities.InventoryEntity, error)
	UpdateInventory(id uint, data interface{}) (entities.InventoryEntity, error)
}

type inventoryService struct {
	inventoryRepository repositories.InventoryRepository
}

func (s *inventoryService) ListInventories() ([]entities.InventoryEntity, error) {
	inventories, appErr := s.inventoryRepository.List()
	if appErr != nil {
		return nil, appErr
	}
	return inventories, nil
}

func (s *inventoryService) GetInventoryByID(id uint) (entities.InventoryEntity, error) {
	// TODO: Implement logic to get inventory by ID
	return entities.InventoryEntity{}, nil
}

func (s *inventoryService) CreateInventory(data interface{}) (entities.InventoryEntity, error) {
	// TODO: Implement logic to create inventory
	return entities.InventoryEntity{}, nil
}

func (s *inventoryService) UpdateInventory(id uint, data interface{}) (entities.InventoryEntity, error) {
	// TODO: Implement logic to update inventory
	return entities.InventoryEntity{}, nil
}

func NewInventoryService(
	repo repositories.InventoryRepository,
) InventoryService {
	return &inventoryService{
		inventoryRepository: repo,
	}
}
