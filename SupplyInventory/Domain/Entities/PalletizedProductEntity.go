package entities

type PalletizedProductEntity struct {
	ID       uint `gorm:"primaryKey"`
	PalletID uint
	EAN      int `json:"ean" binding:"required"`
	Quantity int `json:"quantity" binding:"required"`
}

func (PalletizedProductEntity) TableName() string {
	return "palletized_products"
}
