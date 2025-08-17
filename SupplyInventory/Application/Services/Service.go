package services

import (
	domain "go_inventory/SupplyInventory/Domain"
	infrastructure "go_inventory/SupplyInventory/Infrastructure"
)

func ListPositions() []domain.Position {
	positions, err := infrastructure.GetAllPositions()
	if err != nil {
		return nil
	}

	return positions
}

func FindPositionById(id uint) (*domain.Position, error) {
	position, err := infrastructure.GetSupplyById(id)
	if err != nil {
		return nil, err
	}

	return position, nil
}

func CreatePosition(position domain.Position) (*domain.Position, error) {
	newPosition, err := infrastructure.AddSupply(position)
	if err != nil {
		return nil, err
	}

	return newPosition, nil
}

func AddProductsToPositition(product domain.PositionProduct) *domain.Position {
	newPositionProduct, err := infrastructure.AddProductsToPosition(product)
	if err != nil {
		return nil
	}

	return newPositionProduct
}
