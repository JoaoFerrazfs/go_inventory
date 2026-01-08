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
		req, _ := http.NewRequest("POST", "/api/v1/pallet-racks/", bytes.NewBuffer(body))
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
