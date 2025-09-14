package domain

type PalletEntity struct {
	ID                uint                      `gorm:"primaryKey" json:"id"`
	Name              string                    `gorm:"unique;not null" json:"name" binding:"required"`
	PalletizedProduct []PalletizedProductEntity `gorm:"constraint:OnDelete:CASCADE;foreignKey:PalletID" json:"palletizedProduct"`
	PalletRackID      uint                      `gorm:"not null" json:"palletRackId" binding:"required"`
	QrCode            string                    `json:"qr_code"`
	QrCodeUrl         string                    `json:"qr_code_url"`
}

func (PalletEntity) TableName() string {
	return "pallets"
}
