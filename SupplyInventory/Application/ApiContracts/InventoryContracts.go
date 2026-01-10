package apiContracts

import "time"

// InventoryResponse represents the response structure for inventory-related endpoints.
type InventoryResponse struct {
	Data InventoryData `json:"data"`
}

// InventoryData represents a single inventory object in the response.
type InventoryData struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    string    `json:"status"`
	StartedAt time.Time `json:"started_at"`
	Name      string    `json:"name"`
	UserID    uint      `json:"user_id"`
	User      UserData  `json:"user"`
}

// UserData represents user information in the inventory response.
type UserData struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// InventoryListResponse represents the response structure for listing inventories.
type InventoryListResponse struct {
	Data struct {
		Inventories []InventoryData `json:"inventories"`
	} `json:"data"`
}
