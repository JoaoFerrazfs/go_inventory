package requests

type CreateProductRequest struct {
	EAN  string `json:"ean" binding:"required,len=13"`
	Name string `json:"name" binding:"required"`
}
