package entities

import "time"

type PalletizedProductEntity struct {
	ID        uint `gorm:"primaryKey"`
	PalletID  uint
	InventoryID uint `gorm:"index" json:"inventory_id,omitempty"`
	EAN       int       `json:"ean" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (PalletizedProductEntity) TableName() string {
	return "palletized_products"
}
