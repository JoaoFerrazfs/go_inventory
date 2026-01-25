package apiContracts

// ProductResponse represents the response structure for product-related endpoints.
type ProductResponse struct {
	Data ProductData `json:"data"`
}

// ProductData represents a single product object in the response.
type ProductData struct {
	ID   uint   `json:"id"`
	EAN  string `json:"ean"`
	Name string `json:"name"`
}
