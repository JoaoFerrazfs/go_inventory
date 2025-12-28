package entities

type PalletRackEntity struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"unique;not null" json:"name" binding:"required"`
	Pallets       []PalletEntity `gorm:"foreignKey:PalletRackID" json:"pallets"`
	Location      string         `gorm:"not null" json:"location"  binding:"required" `
	TotalCapacity int            `gorm:"not null;check:total_capacity > 0" json:"total_capacity" binding:"required"`
}

func (PalletRackEntity) TableName() string {
	return "pallet_racks"
}

func NewPalletRackEntity(
	name string,
	location string,
	totalCapacity int,
) *PalletRackEntity {
	return &PalletRackEntity{
		Name:          name,
		Location:      location,
		TotalCapacity: totalCapacity,
	}
}
