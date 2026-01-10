package entities

import "time"

type InventoryStatus string

const (
	InventoryStatusOpen   InventoryStatus = "Open"
	InventoryStatusClosed InventoryStatus = "Closed"
)

type InventoryEntity struct {
	ID          uint               `gorm:"primaryKey" json:"id"`
	PalletRacks []PalletRackEntity `gorm:"foreignKey:InventoryID;constraint:OnDelete:CASCADE;" json:"pallet_racks,omitempty"`
	CreatedAt   time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	Status      InventoryStatus    `gorm:"type:varchar(50);default:'Open'" json:"status"`
	StartedAt   time.Time          `json:"started_at,omitempty"`
	EndedAt     *time.Time         `json:"ended_at,omitempty"`
	Name        string             `gorm:"type:varchar(255);not null" json:"name"`
	Description string             `gorm:"type:text" json:"description,omitempty"`
	UserID      uint               `gorm:"not null" json:"user_id"`
	User        UserEntity         `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

func (InventoryEntity) TableName() string {
	return "inventories"
}
