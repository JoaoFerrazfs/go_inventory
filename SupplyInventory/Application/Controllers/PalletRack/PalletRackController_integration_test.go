//go:build integration
// +build integration

package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	palletRackRequests "go_inventory/SupplyInventory/Application/Requests/PalletRack"
	integration "go_inventory/SupplyInventory/tests/integration"
)

func TestIntegration_CreatePalletRack(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)
	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		r := h.SetupRouterForPalletRack(tx)

		createReq := palletRackRequests.PalletRackRequest{
			Name:          "Rack1",
			Location:      "A1",
			TotalCapacity: 100,
		}
		body, _ := json.Marshal(createReq)

		// Actions
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/racks/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		// create and set inventory header for the request
		inv := h.CreateTestPalletRack(tx, "tmp", "loc", 1)
		req.Header.Set("X-Inventory-ID", fmt.Sprintf("%d", inv.InventoryID))
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusCreated, w.Code)
		return nil
	})
}

func TestIntegration_ListPalletRacks_Filtering(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)
	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		rack1 := h.CreateTestPalletRack(tx, "Rack1", "Loc1", 10)
		h.CreateTestPalletRack(tx, "Rack2", "Loc2", 10) // This creates a NEW inventory internally

		r := h.SetupRouterForPalletRack(tx)

		// Actions
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/racks/", nil)
		req.Header.Set("X-Inventory-ID", fmt.Sprintf("%d", rack1.InventoryID))
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code) 

		var response struct {
			Data  []interface{} `json:"data"`
			Total int64         `json:"total"`
		}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 1, len(response.Data))
		assert.Equal(t, int64(1), response.Total)

		return nil
	})
}
