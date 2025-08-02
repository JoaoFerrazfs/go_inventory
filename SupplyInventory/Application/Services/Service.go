
package services

import (
	"go_inventory/SupplyInventory/Domain"
	"go_inventory/SupplyInventory/Infrastructure"
)

func ListPositions() []domain.Position {
	return infrastructure.GetAllPositions()
}

func FindPositionById(id string) *domain.Position {
	return infrastructure.GetSupplyById(id)
}

func CreatePosition(position domain.Position) {
	infrastructure.AddSupply(position)
}
