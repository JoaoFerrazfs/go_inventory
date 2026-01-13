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

type PaginatedRacksResponse struct {
	Data  []TransformedRack `json:"data"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}
