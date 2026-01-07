package entities

import "time"

type InventoryEntity struct {
	ID          uint               `gorm:"primaryKey" json:"id"`
	PalletRacks []PalletRackEntity `gorm:"foreignKey:id" json:"pallet_racks"`
	CreatedAt   time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
}
