package infrastructure

import (
	domain "go_inventory/SupplyInventory/Domain"
	db "go_inventory/SupplyInventory/Infrastructure/Db"
)

func GetAllPositions() ([]domain.Position, error) {
	var positions []domain.Position
	if err := db.DB.Find(&positions).Error; err != nil {
		return nil, err
	}
	return positions, nil
}

func GetSupplyById(id string) (*domain.Position, error) {
	var position domain.Position
	if err := db.DB.First(&position, id).Error; err != nil {
		return nil, err
	}
	return &position, nil
}

func AddSupply(position domain.Position) (*domain.Position, error) {
	if err := db.DB.Create(&position).Error; err != nil {
		return nil, err
	}
	return &position, nil
}
