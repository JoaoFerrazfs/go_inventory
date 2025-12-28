package apiContracts

import entities "go_inventory/SupplyInventory/Domain/Entities"

type TransformedRack struct {
	ID            uint
	Name          string
	Pallets       []entities.PalletEntity
	Location      string
	TotalCapacity int
	PercetageUsed float64
}
