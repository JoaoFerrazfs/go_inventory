package domain

type Position struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name" binding:"required"`
	Stock int    `json:"stock" binding:"required"`
	EAN   int    `json:"ean" binding:"required"`
}
