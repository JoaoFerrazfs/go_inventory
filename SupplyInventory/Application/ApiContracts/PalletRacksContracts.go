package apiContracts

import domain "go_inventory/SupplyInventory/Domain"

type TransformedRack struct {
	ID            uint
	Name          string
	Pallets       []domain.PalletEntity
	Location      string
	TotalCapacity int
	PercetageUsed float64
}
