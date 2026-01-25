package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	errors "go_inventory/Helpers/Errors"
	palletizedproduct "go_inventory/SupplyInventory/Application/Controllers/PalletizedProduct"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	testutils "go_inventory/SupplyInventory/tests/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- TESTES ADICIONADOS NO FINAL DO ARQUIVO ---

func TestDeleteProductsFromPallet_AppError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletizedProductService)
	controller := palletizedproduct.NewPalletizedProductController(mockService)
	// Expectations
	mockService.On("DeleteProductsFromPallet", uint(1), 123).Return(false, errors.NewAppError("fail", 400))

	// Actions
	r := gin.Default()
	group := r.Group("/")
	controller.RegisterProductPallet(group)
	req, _ := http.NewRequest("DELETE", "/1/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddProductsToPallet_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletizedProductService)
	controller := palletizedproduct.NewPalletizedProductController(mockService)
	// Expectations
	mockService.On("AddProductsToPallet", uint(1), 123, 10, testutils.DefaultInventoryID).Return(nil, errors.NewAppError("fail", 500))

	// Actions
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("inventoryID", testutils.DefaultInventoryID)
		c.Next()
	})
	group := r.Group("/")
	controller.RegisterProductPallet(group)
	body, _ := json.Marshal(map[string]interface{}{"EAN": 123, "Quantity": 10})
	req, _ := http.NewRequest("PATCH", "/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestAddProductsToPallet_InvalidPalletId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletizedProductService)
	controller := palletizedproduct.NewPalletizedProductController(mockService)
	r := gin.Default()
	group := r.Group("/")
	controller.RegisterProductPallet(group)
	body, _ := json.Marshal(map[string]interface{}{"ean": 123, "quantity": 10})
	req, _ := http.NewRequest("PATCH", "/invalid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

// Testa erro de palletId inválido no DELETE
func TestDeleteProductsFromPallet_InvalidPalletId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletizedProductService)
	controller := palletizedproduct.NewPalletizedProductController(mockService)
	r := gin.Default()
	group := r.Group("/")
	controller.RegisterProductPallet(group)
	req, _ := http.NewRequest("DELETE", "/invalid/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Testa erro de productsEan inválido no DELETE
func TestDeleteProductsFromPallet_InvalidProductsEan(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletizedProductService)
	controller := palletizedproduct.NewPalletizedProductController(mockService)
	r := gin.Default()
	group := r.Group("/")
	controller.RegisterProductPallet(group)
	req, _ := http.NewRequest("DELETE", "/1/invalid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

type mockPalletizedProductService struct {
	mock.Mock
}

func (m *mockPalletizedProductService) AddProductsToPallet(palletId uint, ean int, quantity int, inventoryID uint) (*entities.PalletEntity, *errors.AppError) {
	args := m.Called(palletId, ean, quantity, inventoryID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*entities.PalletEntity), args.Get(1).(*errors.AppError)
}

func (m *mockPalletizedProductService) DeleteProductsFromPallet(palletId uint, productsEan int) (bool, *errors.AppError) {
	args := m.Called(palletId, productsEan)
	return args.Bool(0), args.Get(1).(*errors.AppError)
}

func TestAddProductsToPallet_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set
	mockService := new(mockPalletizedProductService)
	controller := palletizedproduct.NewPalletizedProductController(mockService)
	pallet := &entities.PalletEntity{ID: 1, Name: "Test Pallet"}

	// Expectations
	mockService.On("AddProductsToPallet", uint(1), 1234567891234, 10, testutils.DefaultInventoryID).Return(pallet, (*errors.AppError)(nil))

	// Actions
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("inventoryID", testutils.DefaultInventoryID)
		c.Next()
	})
	group := r.Group("/")
	controller.RegisterProductPallet(group)
	body, _ := json.Marshal(map[string]interface{}{"EAN": 1234567891234, "Quantity": 10})
	req, _ := http.NewRequest("PATCH", "/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
}
