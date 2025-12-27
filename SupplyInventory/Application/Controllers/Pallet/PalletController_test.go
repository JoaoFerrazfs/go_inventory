package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	errors "go_inventory/Helpers/Errors"
	pallet "go_inventory/SupplyInventory/Application/Controllers/Pallet"
	domain "go_inventory/SupplyInventory/Domain"
)

type mockPalletService struct {
	mock.Mock
}

func (m *mockPalletService) ListPallets() ([]domain.PalletEntity, *errors.AppError) {
	args := m.Called()
	return args.Get(0).([]domain.PalletEntity), args.Get(1).(*errors.AppError)
}
func (m *mockPalletService) FindPalletById(id uint) (*domain.PalletEntity, *errors.AppError) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*domain.PalletEntity), args.Get(1).(*errors.AppError)
}
func (m *mockPalletService) DeletePalletById(id uint) (bool, *errors.AppError) {
	args := m.Called(id)
	return args.Bool(0), args.Get(1).(*errors.AppError)
}
func (m *mockPalletService) CreatePallet(name string, rackID uint) (*domain.PalletEntity, *errors.AppError) {
	args := m.Called(name, rackID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*domain.PalletEntity), args.Get(1).(*errors.AppError)
}
func (m *mockPalletService) UpdatePallet(id uint, name string, rackID uint) (*domain.PalletEntity, *errors.AppError) {
	args := m.Called(id, name, rackID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*domain.PalletEntity), args.Get(1).(*errors.AppError)
}

func TestListPallets_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set
	mockService := new(mockPalletService)
	controller := pallet.NewPalletController(mockService)
	pallets := []domain.PalletEntity{{ID: 1, Name: "Pallet1"}}

	// Expectations
	mockService.On("ListPallets").Return(pallets, (*errors.AppError)(nil))

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	controller.ListPallets(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
}
