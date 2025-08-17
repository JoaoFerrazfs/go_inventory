package domain

type PositionProduct struct {
	ID         uint `gorm:"primaryKey"`
	PositionID uint
	EAN        int `json:"ean" binding:"required"`
	Quantity   int `json:"quantity" binding:"required"`
}
