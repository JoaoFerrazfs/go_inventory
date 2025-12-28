package pallet

type PalletRequest struct {
	Name         string `json:"name" binding:"required"`
	PalletRackID uint   `json:"palletRackId" binding:"required"`
}
