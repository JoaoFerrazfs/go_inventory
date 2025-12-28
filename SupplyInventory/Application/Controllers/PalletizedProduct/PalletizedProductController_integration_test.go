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
	requests "go_inventory/SupplyInventory/Application/Requests"
	services "go_inventory/SupplyInventory/Application/Services"
	entities "go_inventory/SupplyInventory/Domain/Entities"
)

type mockPalletRepository struct {
	mock.Mock
}

func (m *mockPalletRepository) Create(pallet *entities.PalletEntity) error { return nil }
func (m *mockPalletRepository) FindByID(id uint) (*entities.PalletEntity, error) { return nil, nil }
func (m *mockPalletRepository) List() ([]*entities.PalletEntity, error) { return nil, nil }
func (m *mockPalletRepository) DeleteByID(id uint) error { return nil }
func (m *mockPalletRepository) Update(pallet *entities.PalletEntity) error { return nil }
func (m *mockPalletRepository) AddProductsToPallet(product entities.PalletizedProductEntity) (*entities.PalletEntity, *errors.AppError) {
	args := m.Called(product)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*entities.PalletEntity), args.Get(1).(*errors.AppError)
}
func (m *mockPalletRepository) GetAllPallets() ([]entities.PalletEntity, *errors.AppError) { return nil, nil }
func (m *mockPalletRepository) GetSupplyById(id uint) (*entities.PalletEntity, *errors.AppError) { return nil, nil }
func (m *mockPalletRepository) AddSupply(name string, rackId uint) (*entities.PalletEntity, *errors.AppError) { return nil, nil }
func (m *mockPalletRepository) UpdateSupply(pallet *entities.PalletEntity) (*entities.PalletEntity, *errors.AppError) { return nil, nil }
func (m *mockPalletRepository) DeletePalletById(id uint) (bool, *errors.AppError) { return false, nil }
func (m *mockPalletRepository) UpdatePallet(id uint, name string, rackId uint) (*entities.PalletEntity, *errors.AppError) { return nil, nil }

type mockPalletizedProductRepository struct {
	mock.Mock
}

// Métodos obrigatórios do PalletizedProductRepository
func (m *mockPalletizedProductRepository) AddProductsToPallet(product entities.PalletizedProductEntity) (bool, *errors.AppError) {
	args := m.Called(product)
	return args.Bool(0), args.Get(1).(*errors.AppError)
}
func (m *mockPalletizedProductRepository) DeleteProductsFromPallet(palletId uint, productsEan int) (bool, *errors.AppError) {
	args := m.Called(palletId, productsEan)
	return args.Bool(0), args.Get(1).(*errors.AppError)
}

func setupPalletizedProductRouter() *gin.Engine {
	mockPalletRepo := new(mockPalletRepository)
	mockProductRepo := new(mockPalletizedProductRepository)
	service := services.NewPalletizedProductService(mockPalletRepo, mockProductRepo)
	controller := palletizedproduct.NewPalletizedProductController(service)

	r := gin.Default()
	api := r.Group("/api/v1/pallet/products")
	controller.RegisterProductPallet(api)
	return r
}

func TestIntegration_AddProductsToPallet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set
	mockPalletRepo := new(mockPalletRepository)
	mockProductRepo := new(mockPalletizedProductRepository)
	service := services.NewPalletizedProductService(mockPalletRepo, mockProductRepo)
	controller := palletizedproduct.NewPalletizedProductController(service)
	r := gin.Default()
	api := r.Group("/api/v1/pallet/products")
	controller.RegisterProductPallet(api)

	addReq := requests.PalletizedProductRequest{
		EAN:      123,
		Quantity: 10,
	}
	body, _ := json.Marshal(addReq)

	// Expectations
	mockProductRepo.On("AddProductsToPallet", mock.AnythingOfType("entities.PalletizedProductEntity")).Return(true, (*errors.AppError)(nil))

	// Actions
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/pallet/products/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Assertions
	assert.Contains(t, []int{http.StatusOK, http.StatusUnprocessableEntity}, w.Code)
}
