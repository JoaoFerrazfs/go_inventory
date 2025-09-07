package requests

type PalletizedProductRequest struct {
	EAN      int `json:"ean" binding:"required"`
	Quantity int `json:"quantity" binding:"required"`
}
