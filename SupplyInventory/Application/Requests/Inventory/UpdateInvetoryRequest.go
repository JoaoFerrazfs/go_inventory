package inventory

type UpdateInventoryRequest struct {
	Name        string `json:"name" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
	Status      string `json:"status" binding:"omitempty,oneof=Open Closed"`
}
