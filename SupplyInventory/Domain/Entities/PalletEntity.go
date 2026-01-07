package entities

import "time"

type PalletEntity struct {
	ID                uint                      `gorm:"primaryKey" json:"id"`
	Name              string                    `gorm:"unique;not null" json:"name" binding:"required"`
	PalletizedProduct []PalletizedProductEntity `gorm:"constraint:OnDelete:CASCADE;foreignKey:PalletID" json:"palletizedProduct"`
	PalletRackID      uint                      `gorm:"not null" json:"palletRackId" binding:"required"`
	PalletRackName    string                    `gorm:"not null" json:"palletRackName" binding:"required"`
	QrCode            string                    `json:"qr_code"`
	QrCodeUrl         string                    `json:"qr_code_url"`
	CreatedAt         time.Time                 `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time                 `gorm:"autoUpdateTime" json:"updated_at"`
}

func (PalletEntity) TableName() string {
	return "pallets"
}
