//go:build integration
// +build integration

package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	inventoryRequests "go_inventory/SupplyInventory/Application/Requests/Inventory"
	integration "go_inventory/SupplyInventory/tests/integration"
)

func TestIntegration_ListInventories(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	defer h.Stop()
	h.TruncateTables(h.DB)

	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		h.CreateTestInventory(tx)
		r := h.SetupRouterForInventory(tx)

		// Actions
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/inventories/", nil)
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		data := response["data"].(map[string]interface{})
		inventories := data["inventories"].([]interface{})
		assert.GreaterOrEqual(t, len(inventories), 1)

		return nil
	})
}

func TestIntegration_GetInventoryById(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	defer h.Stop()
	h.TruncateTables(h.DB)

	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		inv := h.CreateTestInventory(tx)
		r := h.SetupRouterForInventory(tx)

		// Actions
		w := httptest.NewRecorder()
		url := fmt.Sprintf("/api/v1/inventories/%d", inv.ID)
		req, _ := http.NewRequest("GET", url, nil)
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		data := response["data"].(map[string]interface{})
		inventory := data["inventory"].(map[string]interface{})
		assert.Equal(t, float64(inv.ID), inventory["id"])
		assert.Equal(t, inv.Name, inventory["name"])

		return nil
	})
}

func TestIntegration_CreateInventory(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	defer h.Stop()
	h.TruncateTables(h.DB)

	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		user := h.CreateTestUser(tx, "Test User", "test@example.com", "password")
		r := h.SetupRouterForInventory(tx)

		createReq := inventoryRequests.InventoryRequest{
			Name:        "New Inventory",
			Description: "Inventory Description",
		}
		body, _ := json.Marshal(createReq)

		// Actions
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/inventories/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Test-User-ID", strconv.Itoa(int(user.ID)))
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		data := response["data"].(map[string]interface{})
		inventory := data["inventory"].(map[string]interface{})
		assert.Equal(t, "New Inventory", inventory["name"])

		return nil
	})
}

func TestIntegration_UpdateInventory(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	defer h.Stop()
	h.TruncateTables(h.DB)

	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		inv := h.CreateTestInventory(tx)
		r := h.SetupRouterForInventory(tx)

		// Actions
		w := httptest.NewRecorder()
		url := fmt.Sprintf("/api/v1/inventories/%d", inv.ID)
		req, _ := http.NewRequest("PUT", url, nil)
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Update inventory")

		return nil
	})
}
