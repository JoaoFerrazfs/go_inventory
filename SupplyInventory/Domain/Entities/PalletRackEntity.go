package entities

import "time"

type PalletRackEntity struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	InventoryID   uint           `gorm:"index;not null" json:"inventory_id" binding:"required"`
	Name          string         `gorm:"unique;not null" json:"name" binding:"required"`
	Pallets       []PalletEntity `gorm:"foreignKey:PalletRackID" json:"pallets"`
	Location      string         `gorm:"not null" json:"location"  binding:"required" `
	TotalCapacity int            `gorm:"not null;check:total_capacity > 0" json:"total_capacity" binding:"required"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

func (PalletRackEntity) TableName() string {
	return "pallet_racks"
}

func NewPalletRackEntity(inventoryID uint,
	name string,
	location string,
	totalCapacity int,
) *PalletRackEntity {
	return &PalletRackEntity{
		InventoryID:   inventoryID,
		Name:          name,
		Location:      location,
		TotalCapacity: totalCapacity,
	}
}
