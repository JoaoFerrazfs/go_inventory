package domain

type Pallet struct {
	ID                uint                      `gorm:"primaryKey" json:"id"`
	Name              string                    `gorm:"unique;not null" json:"name" binding:"required"`
	PalletizedProduct []PalletizedProductEntity `gorm:"foreignKey:PalletID" json:"palletizedProduct"`
	QrCode            string                    `json:"qr_code"`
	QrCodeUrl         string                    `json:"qr_code_url"`
}
