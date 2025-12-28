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

	userRequests "go_inventory/SupplyInventory/Application/Requests/User"
	integration "go_inventory/SupplyInventory/tests/integration"
)

func TestIntegration_CreateUser(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)
	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		r := h.SetupRouterForUser(tx)

		createReq := userRequests.UserRequest{
			Name:     "Admin",
			Email:    "admin@example.com",
			Password: "admin123",
		}
		body, _ := json.Marshal(createReq)

		// Actions
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/create", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusCreated, w.Code)
		return nil
	})
}
