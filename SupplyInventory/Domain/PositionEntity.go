package domain

type Position struct {
	ID        uint              `gorm:"primaryKey" json:"id"`
	Name      string            `gorm:"unique;not null" json:"name" binding:"required"`
	Products  []PositionProduct `gorm:"foreignKey:PositionID" json:"products"`
	QrCode    string            `json:"qr_code"`
	QrCodeUrl string            `json:"qr_code_url"`
}
