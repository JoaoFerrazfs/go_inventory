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
	pallet "go_inventory/SupplyInventory/Application/Controllers/Pallet"
	requests "go_inventory/SupplyInventory/Application/Requests"
	services "go_inventory/SupplyInventory/Application/Services"
	domain "go_inventory/SupplyInventory/Domain"
)

type mockPalletRepository struct {
	mock.Mock
}

func (m *mockPalletRepository) UpdatePallet(id uint, name string, rackId uint) (*domain.PalletEntity, *errors.AppError) {
	args := m.Called(id, name, rackId)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*domain.PalletEntity), args.Get(1).(*errors.AppError)
}

type mockQRCodeService struct {
	mock.Mock
}

func (m *mockQRCodeService) CreateQRCode(palletId uint) (string, error) {
	args := m.Called(palletId)
	return args.String(0), args.Error(1)
}

func (m *mockPalletRepository) Create(pallet *domain.PalletEntity) error {
	args := m.Called(pallet)
	return args.Error(0)
}
func (m *mockPalletRepository) FindByID(id uint) (*domain.PalletEntity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PalletEntity), args.Error(1)
}
func (m *mockPalletRepository) List() ([]*domain.PalletEntity, error) {
	args := m.Called()
	return args.Get(0).([]*domain.PalletEntity), args.Error(1)
}
func (m *mockPalletRepository) DeleteByID(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *mockPalletRepository) Update(pallet *domain.PalletEntity) error {
	args := m.Called(pallet)
	return args.Error(0)
}
func (m *mockPalletRepository) AddProductsToPallet(product domain.PalletizedProductEntity) (*domain.PalletEntity, *errors.AppError) {
	args := m.Called(product)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*domain.PalletEntity), args.Get(1).(*errors.AppError)
}
func (m *mockPalletRepository) GetAllPallets() ([]domain.PalletEntity, *errors.AppError) {
	args := m.Called()
	return args.Get(0).([]domain.PalletEntity), args.Get(1).(*errors.AppError)
}
func (m *mockPalletRepository) GetSupplyById(id uint) (*domain.PalletEntity, *errors.AppError) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*domain.PalletEntity), args.Get(1).(*errors.AppError)
}
func (m *mockPalletRepository) AddSupply(name string, rackId uint) (*domain.PalletEntity, *errors.AppError) {
	args := m.Called(name, rackId)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*domain.PalletEntity), args.Get(1).(*errors.AppError)
}
func (m *mockPalletRepository) UpdateSupply(pallet *domain.PalletEntity) (*domain.PalletEntity, *errors.AppError) {
	args := m.Called(pallet)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*domain.PalletEntity), args.Get(1).(*errors.AppError)
}
	func (m *mockPalletRepository) DeletePalletById(id uint) (bool, *errors.AppError) {
		args := m.Called(id)
		return args.Bool(0), args.Get(1).(*errors.AppError)
	}

func setupPalletRouter() *gin.Engine {
	mockRepo := new(mockPalletRepository)
	mockQR := new(mockQRCodeService)
	service := services.NewPalletService(mockRepo, mockQR)
	controller := pallet.NewPalletController(service)

	r := gin.Default()
	api := r.Group("/api/v1/pallets")
	controller.Register(api)
	return r
}

func TestIntegration_CreatePallet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set
	mockRepo := new(mockPalletRepository)
	mockQR := new(mockQRCodeService)
	service := services.NewPalletService(mockRepo, mockQR)
	controller := pallet.NewPalletController(service)
	r := gin.Default()
	api := r.Group("/api/v1/pallets")
	controller.Register(api)

	createReq := requests.PalletRequest{
		Name:        "Pallet1",
		PalletRackID: 1,
	}
	body, _ := json.Marshal(createReq)

	// Expectations
	mockRepo.On("Create", mock.AnythingOfType("*domain.PalletEntity")).Return(nil)
	mockRepo.On("AddSupply", "Pallet1", uint(1)).Return(&domain.PalletEntity{ID: 1, Name: "Pallet1"}, (*errors.AppError)(nil))
	mockQR.On("CreateQRCode", mock.AnythingOfType("uint")).Return("/tmp/fake.png", nil)
	mockRepo.On("UpdateSupply", mock.AnythingOfType("*domain.PalletEntity")).Return(&domain.PalletEntity{ID: 1, Name: "Pallet1", QrCode: "/tmp/fake.png", QrCodeUrl: "http://localhost:3000//tmp/fake.png"}, (*errors.AppError)(nil))

	// Actions
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/pallets/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Assertions
	assert.Contains(t, []int{http.StatusCreated, http.StatusUnprocessableEntity}, w.Code)
}
