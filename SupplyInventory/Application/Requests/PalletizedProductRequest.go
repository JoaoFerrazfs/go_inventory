package requests

type PalletizedProductRequest struct {
	PalletID uint `json:"palletId" binding:"required"`
	EAN      int  `json:"ean" binding:"required"`
	Quantity int  `json:"quantity" binding:"required"`
}
