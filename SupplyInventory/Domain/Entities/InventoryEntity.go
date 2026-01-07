package entities

import "time"

type InventoryEntity struct {
	ID          uint               `gorm:"primaryKey" json:"id"`
	PalletRacks []PalletRackEntity `gorm:"foreignKey:InventoryID;constraint:OnDelete:CASCADE;" json:"pallet_racks,omitempty"`
	CreatedAt   time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time          `gorm:"autoUpdateTime" json:"updated_at"`

	Status    string     `gorm:"size:50;default:'open'" json:"status"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
}

func (InventoryEntity) TableName() string {
	return "inventories"
}
