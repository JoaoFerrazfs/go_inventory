package infrastructure

import (
	"go_inventory/SupplyInventory/Domain"
)

var positition = []domain.Position{
	{
		ID:    "1",
		Name:  "Position A",
		Stock: 100,
	},
	{
		ID:    "2",
		Name:  "Position B",
		Stock: 100,
	},
}

func GetAllPositions() []domain.Position {
	return positition
}

func GetSupplyById(id string) *domain.Position {
	for _, position := range positition {
		if position.ID == id {
			return &position
		}
}
	return nil
}

func AddSupply(position domain.Position) {
	positition = append(positition, position)
}