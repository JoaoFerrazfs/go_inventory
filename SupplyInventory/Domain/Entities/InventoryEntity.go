package entities

import "time"

type InventoryEntity struct {
	ID          uint               `gorm:"primaryKey" json:"id"`
	PalletRacks []PalletRackEntity `gorm:"foreignKey:InventoryID;constraint:OnDelete:CASCADE;" json:"pallet_racks,omitempty"`
	CreatedAt   time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
<<<<<<< Updated upstream

	Status    string     `gorm:"size:50;default:'open'" json:"status"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
=======
	Status      InventoryStatus    `gorm:"type:varchar(50);default:'Open'" json:"status"`
	StartedAt   time.Time          `json:"started_at,omitempty"`
	EndedAt     *time.Time         `json:"ended_at,omitempty"`
	Name        string             `gorm:"type:varchar(255);not null" json:"name"`
	Description string             `gorm:"type:text" json:"description,omitempty"`
>>>>>>> Stashed changes
}

func (InventoryEntity) TableName() string {
	return "inventories"
}
