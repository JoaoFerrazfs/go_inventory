//go:build integration
// +build integration

package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	palletRequests "go_inventory/SupplyInventory/Application/Requests/Pallet"
	integration "go_inventory/SupplyInventory/tests/integration"
)

func TestIntegration_CreatePallet(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)
	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		rack := h.CreateTestPalletRack(tx, "Rack1_Pallet", "Location1", 10)
		r := h.SetupRouterForPallet(tx)

		createReq := palletRequests.PalletRequest{
			Name:         "Pallet1",
			PalletRackID: rack.ID,
		}
		body, _ := json.Marshal(createReq)

		// Actions
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/pallets/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusCreated, w.Code)
		return nil
	})
}

func TestIntegration_ExportPalletsCsv(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)
	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		rack := h.CreateTestPalletRack(tx, "Rack1_Export", "Location1", 10)
		h.CreateTestPallet(tx, "Pallet1_Export", rack.ID)
		r := h.SetupRouterForPallet(tx)

		// Actions
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/pallets/export", nil)
		r.ServeHTTP(w, req)

		// Assertions
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].(map[string]interface{})
		url := data["url"].(string)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, url, "http://localhost:3000/reports/Pallets")
		assert.Contains(t, url, ".csv")

		return nil
	})
}
