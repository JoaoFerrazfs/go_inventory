package requests

type PalletRackRequest struct {
	Name          string `binding:"required"`
	Location      string `binding:"required" `
	TotalCapacity int    `binding:"required,min=1"`
}
