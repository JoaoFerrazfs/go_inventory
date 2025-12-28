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

	authRequests "go_inventory/SupplyInventory/Application/Requests/Auth"
	authResponses "go_inventory/SupplyInventory/Application/Responses/Auth"
	integration "go_inventory/SupplyInventory/tests/integration"
)

func TestIntegration_Login(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)
	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		testUser := h.CreateTestUser(tx, "Admin", "admin@example.com", "admin123")
		r := h.SetupRouterForAuth(tx)

		loginReq := authRequests.LoginRequest{
			Email:    testUser.Email,
			Password: "admin123",
		}

		body, _ := json.Marshal(loginReq)

		// Actions
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		var resp authResponses.AuthResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NotEmpty(t, resp.Token)
		return nil
	})
}

func TestIntegration_RefreshToken(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		h.TruncateTables(tx)
		r := h.SetupRouterForAuth(tx)

		refreshReq := map[string]string{"refreshToken": "sometoken"}
		body, _ := json.Marshal(refreshReq)

		// Actions
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/refreshToken", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		// Assertions
		assert.Contains(t, []int{http.StatusOK, http.StatusUnprocessableEntity, http.StatusUnauthorized}, w.Code)
		return nil
	})
}
