//go:build integration
// +build integration

package controllers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	integration "go_inventory/SupplyInventory/tests/integration"
)

func TestAdminPalletRackController_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	h := integration.NewIntegrationTestHelper()
	defer h.Stop()
	h.TruncateTables(h.DB)

	t.Run("list racks without inventory_id", func(t *testing.T) {
		h.DB.Transaction(func(tx *gorm.DB) error {
			// Set
			// Create test data
			h.CreateTestPalletRack(tx, "Rack1", "Loc1", 10)
			h.CreateTestPalletRack(tx, "Rack2", "Loc2", 10)

			r := h.SetupRouterForAdminPalletRack(tx)

			// Actions
			req, _ := http.NewRequest("GET", "/api/v1/admin/racks/", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, http.StatusOK, w.Code)
			// Assert response body contains racks
			return nil
		})
	})

	t.Run("list racks with inventory_id", func(t *testing.T) {
		h.DB.Transaction(func(tx *gorm.DB) error {
			// Set
			rack := h.CreateTestPalletRack(tx, "Rack3", "Loc3", 10)

			r := h.SetupRouterForAdminPalletRack(tx)

			// Actions
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/admin/racks/?inventory_id=%d&page=1&limit=10", rack.InventoryID), nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, http.StatusOK, w.Code)
			// Assert pagination and filtered data
			return nil
		})
	})
}
