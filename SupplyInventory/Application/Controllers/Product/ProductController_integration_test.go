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

	productRequests "go_inventory/SupplyInventory/Application/Requests/Product"
	integration "go_inventory/SupplyInventory/tests/integration"
)

func TestIntegration_CreateProduct_Success(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)
	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		r := h.SetupRouterForProduct(tx)

		createReq := productRequests.CreateProductRequest{
			EAN:  "1234567890123",
			Name: "Test Product",
		}
		body, _ := json.Marshal(createReq)

		// Actions
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Test Product", response["data"].(map[string]interface{})["product"].(map[string]interface{})["name"])
		assert.Equal(t, "1234567890123", response["data"].(map[string]interface{})["product"].(map[string]interface{})["ean"])
		return nil
	})
}

func TestIntegration_CreateProduct_DuplicateEAN(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)
	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		r := h.SetupRouterForProduct(tx)

		// Create first product
		h.CreateTestProduct(tx, "1234567890123", "Existing Product")

		createReq := productRequests.CreateProductRequest{
			EAN:  "1234567890123",
			Name: "Duplicate Product",
		}
		body, _ := json.Marshal(createReq)

		// Actions
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Product with given EAN already exists", response["error"])
		return nil
	})
}

func TestIntegration_CreateProduct_InvalidEAN(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)
	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		r := h.SetupRouterForProduct(tx)

		createReq := productRequests.CreateProductRequest{
			EAN:  "123", // Invalid length
			Name: "Test Product",
		}
		body, _ := json.Marshal(createReq)

		// Actions
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		return nil
	})
}

func TestIntegration_GetProductByEAN_Success(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)
	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		r := h.SetupRouterForProduct(tx)

		// Create product
		product := h.CreateTestProduct(tx, "1234567890123", "Test Product")

		// Actions
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/products/1234567890123", nil)
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, product.Name, response["data"].(map[string]interface{})["product"].(map[string]interface{})["name"])
		assert.Equal(t, product.EAN, response["data"].(map[string]interface{})["product"].(map[string]interface{})["ean"])
		return nil
	})
}

func TestIntegration_GetProductByEAN_NotFound(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)
	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		r := h.SetupRouterForProduct(tx)

		// Actions
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/products/1234567890123", nil)
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "not found")
		return nil
	})
}

func TestIntegration_DeleteProduct_Success(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)
	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		r := h.SetupRouterForProduct(tx)

		// Create product
		h.CreateTestProduct(tx, "1234567890123", "Test Product")

		// Actions
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/products/1234567890123", nil)
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Product deleted successfully", response["data"])
		return nil
	})
}

func TestIntegration_DeleteProduct_NotFound(t *testing.T) {
	h := integration.NewIntegrationTestHelper()
	h.TruncateTables(h.DB)
	h.DB.Transaction(func(tx *gorm.DB) error {
		// Set
		r := h.SetupRouterForProduct(tx)

		// Actions
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/products/1234567890123", nil)
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Product not found", response["error"])
		return nil
	})
}
