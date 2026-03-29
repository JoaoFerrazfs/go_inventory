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

	testutils "go_inventory/SupplyInventory/tests/testutils"

	pallet "go_inventory/SupplyInventory/Application/Controllers/Pallet"
	middlewares "go_inventory/SupplyInventory/Application/Middlewares"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"

	palletRequests "go_inventory/SupplyInventory/Application/Requests/Pallet"
	integration "go_inventory/SupplyInventory/tests/integration"
)

func TestIntegration_CreatePallet(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)
	// Set
	rack := h.CreateTestPalletRack(h.DB, "Rack1_Pallet", "Location1", 10)
	r := h.SetupRouterForPallet(h.DB)

	createReq := palletRequests.PalletRequest{
		Name:         "Pallet1",
		PalletRackID: rack.ID,
	}
	body, _ := json.Marshal(createReq)

	// Actions
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/pallets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Inventory-ID", strconv.Itoa(int(rack.InventoryID)))
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_CreatePallet_DifferentInventory(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)

	// Cria dois inventários persistentes
	inv1 := h.CreateTestInventory(h.DB)
	inv2 := h.CreateTestInventory(h.DB)
	// Cria um rack explicitamente atrelado ao inv1
	rack := testutils.CreateTestPalletRack(h.DB, "Rack1", "Location1", 10)
	rack.InventoryID = inv1.ID
	h.DB.Save(rack)

	// Cria router manualmente, sem middleware de autenticação
	controller := pallet.NewPalletController(h.PalletService)
	r := gin.New()
	// Registrar rota diretamente na raiz, sem grupo /api/v1
	invMiddleware := middlewares.NewInventoryMiddleware(h.InventoryRepo)
	r.Use(invMiddleware.Handler())
	controller.Register(r.Group("/pallets"))

	createReq := palletRequests.PalletRequest{
		Name:         "Pallet1",
		PalletRackID: rack.ID,
	}
	body, _ := json.Marshal(createReq)

	// Actions
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/pallets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Inventory-ID", strconv.Itoa(int(inv2.ID))) // Different inventory
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 422, w.Code)
	assert.Contains(t, w.Body.String(), "Pallet Rack does not belong to the same inventory")
}

func TestIntegration_ExportPalletsCsv(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)
	// Set
	rack := h.CreateTestPalletRack(h.DB, "Rack1_Export", "Location1", 10)
	h.CreateTestPallet(h.DB, "Pallet1_Export", rack.ID)
	r := h.SetupRouterForPallet(h.DB)

	// Actions
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/pallets/export", nil)
	req.Header.Set("X-Inventory-ID", strconv.Itoa(int(rack.InventoryID)))
	r.ServeHTTP(w, req)

	// Assertions
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	url := data["url"].(string)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, url, "http://localhost:3000/reports/Pallets")
	assert.Contains(t, url, ".csv")
}

func TestIntegration_UpdatePallet_DifferentInventory(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)
	// Set
	rack1 := h.CreateTestPalletRack(h.DB, "Rack1", "Location1", 10)
	pallet := h.CreateTestPallet(h.DB, "Pallet1", rack1.ID)

	// Create a second inventory
	inv2 := h.CreateTestInventory(h.DB)

	r := h.SetupRouterForPallet(h.DB)

	updateReq := palletRequests.PalletRequest{
		Name:         "Updated Pallet",
		PalletRackID: rack1.ID,
	}
	body, _ := json.Marshal(updateReq)

	// Actions
	w := httptest.NewRecorder()
	url := "/api/v1/pallets/" + strconv.Itoa(int(pallet.ID))
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Inventory-ID", strconv.Itoa(int(inv2.ID))) // Different inventory
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 422, w.Code)
	assert.Contains(t, w.Body.String(), "Pallet Rack does not belong to the same inventory as the pallet")
}
