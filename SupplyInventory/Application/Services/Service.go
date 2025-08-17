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

func FindPositionById(id uint) *domain.Position {
	position, err := infrastructure.GetSupplyById(id)
	if err != nil {
		return nil
	}

	return position
}

func CreatePosition(position domain.Position) *domain.Position {
	newPosition, err := infrastructure.AddSupply(position)
	if err != nil {
		return nil
	}

	return newPosition
}

func AddProductsToPositition(product domain.PositionProduct) *domain.Position {
	newPositionProduct, err := infrastructure.AddProductsToPosition(product)
	if err != nil {
		return nil
	}

	return newPositionProduct
}
