package domain

type Position struct {
	ID    int    `json:"id"`
	Name  string `json:"name" binding:"required"`
	Stock int    `json:"stock" binding:"required"`
	EAN   int    `json:"ean" binding:"required"`
}
