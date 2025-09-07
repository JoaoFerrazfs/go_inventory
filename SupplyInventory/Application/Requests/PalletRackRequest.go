package requests

type PalletRackRequest struct {
	Name string `json:"name" binding:"required"`
}
