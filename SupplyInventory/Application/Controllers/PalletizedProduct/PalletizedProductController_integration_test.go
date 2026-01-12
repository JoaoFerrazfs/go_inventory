//go:build integration
// +build integration

package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	palletizedProductRequests "go_inventory/SupplyInventory/Application/Requests/PalletizedProduct"
	integration "go_inventory/SupplyInventory/tests/integration"
)

func TestIntegration_AddProductsToPallet(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)
	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		rack := h.CreateTestPalletRack(tx, "Rack1_Product", "Location1", 10)
		pallet := h.CreateTestPallet(tx, "Pallet1_Product", rack.ID)
		r := h.SetupRouterForPalletizedProduct(tx)

		addReq := palletizedProductRequests.PalletizedProductRequest{
			EAN:      123,
			Quantity: 10,
		}
		body, _ := json.Marshal(addReq)

		// Actions
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/v1/pallet/products/"+strconv.Itoa(int(pallet.ID)), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Inventory-ID", strconv.Itoa(int(pallet.InventoryID)))
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		return nil
	})
}

func TestIntegration_AddProductsToPallet_DifferentInventory(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)
	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		rack := h.CreateTestPalletRack(tx, "Rack1_Product", "Location1", 10)
		pallet := h.CreateTestPallet(tx, "Pallet1_Product", rack.ID)

		// Create a second inventory
		inv2 := h.CreateTestInventory(tx)

		r := h.SetupRouterForPalletizedProduct(tx)

		addReq := palletizedProductRequests.PalletizedProductRequest{
			EAN:      123,
			Quantity: 10,
		}
		body, _ := json.Marshal(addReq)

		// Actions
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/v1/pallet/products/"+strconv.Itoa(int(pallet.ID)), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Inventory-ID", strconv.Itoa(int(inv2.ID))) // Different inventory
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, 422, w.Code)
		assert.Contains(t, w.Body.String(), "Pallet does not belong to the same inventory")
		return nil
	})
}
