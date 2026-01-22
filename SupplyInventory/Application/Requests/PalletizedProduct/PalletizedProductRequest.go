package palletizedproduct

type PalletizedProductRequest struct {
	EAN      string `json:"ean" binding:"required,len=13"`
	Quantity int    `json:"quantity" binding:"required"`
}
