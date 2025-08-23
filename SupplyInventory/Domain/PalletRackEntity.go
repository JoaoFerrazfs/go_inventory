package domain

type PalletRackEntity struct {
	ID      uint           `gorm:"primaryKey" json:"id"`
	Name    string         `gorm:"unique;not null" json:"name" binding:"required"`
	Pallets []PalletEntity `gorm:"foreignKey:PalletRackID" json:"pallets"`
}

func (PalletRackEntity) TableName() string {
	return "pallet_racks"
}
