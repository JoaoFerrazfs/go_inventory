package palletizedproduct

type PalletizedProductRequest struct {
	EAN      int `json:"ean" binding:"required,min=1000000000000,max=9999999999999"`
	Quantity int `json:"quantity" binding:"required"`
}
