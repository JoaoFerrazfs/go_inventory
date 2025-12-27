package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	errors "go_inventory/Helpers/Errors"
	palletizedproduct "go_inventory/SupplyInventory/Application/Controllers/PalletizedProduct"
	domain "go_inventory/SupplyInventory/Domain"
)

type mockPalletizedProductService struct {
	mock.Mock
}

func (m *mockPalletizedProductService) AddProductsToPallet(palletId uint, ean int, quantity int) (*domain.PalletEntity, *errors.AppError) {
	args := m.Called(palletId, ean, quantity)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*domain.PalletEntity), args.Get(1).(*errors.AppError)
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
	pallet := &domain.PalletEntity{ID: 1, Name: "Test Pallet"}

	// Expectations
	mockService.On("AddProductsToPallet", uint(1), 123, 10).Return(pallet, (*errors.AppError)(nil))

	// Actions
	r := gin.Default()
	group := r.Group("/")
	controller.RegisterProductPallet(group)
	body, _ := json.Marshal(map[string]interface{}{"EAN": 123, "Quantity": 10})
	req, _ := http.NewRequest("PATCH", "/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
}
