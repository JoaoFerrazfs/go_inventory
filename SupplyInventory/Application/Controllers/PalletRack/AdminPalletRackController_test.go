//go:build unit

package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	appErrors "go_inventory/Helpers/Errors"
	apiContracts "go_inventory/SupplyInventory/Application/ApiContracts"
	controllers "go_inventory/SupplyInventory/Application/Controllers/PalletRack"
	entities "go_inventory/SupplyInventory/Domain/Entities"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock para PalletRackService
type MockPalletRackService struct {
	mock.Mock
}

func (m *MockPalletRackService) Create(name string, location string, totalCapacity int, inventoryID uint) (*entities.PalletRackEntity, error) {
	args := m.Called(name, location, totalCapacity, inventoryID)
	return args.Get(0).(*entities.PalletRackEntity), args.Error(1)
}

func (m *MockPalletRackService) ListRacks(inventoryID *uint, page int, limit int) (*apiContracts.PaginatedRacksResponse, error) {
	args := m.Called(inventoryID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*apiContracts.PaginatedRacksResponse), args.Error(1)
}

func (m *MockPalletRackService) FindPalletById(id uint) (*entities.PalletRackEntity, *appErrors.AppError) {
	args := m.Called(id)
	return args.Get(0).(*entities.PalletRackEntity), args.Get(1).(*appErrors.AppError)
}

func (m *MockPalletRackService) DeleteRack(id uint) (bool, *appErrors.AppError) {
	args := m.Called(id)
	return args.Bool(0), args.Get(1).(*appErrors.AppError)
}

func TestAdminPalletRackController_listRacks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success without inventory_id", func(t *testing.T) {
		// Set
		mockService := new(MockPalletRackService)
		controller := controllers.NewAdminPalletRackController(mockService)

		response := &apiContracts.PaginatedRacksResponse{
			Data:  []apiContracts.TransformedRack{},
			Total: 0,
			Page:  1,
			Limit: 10,
		}

		// Expectations
		mockService.On("ListRacks", (*uint)(nil), 1, 10).Return(response, nil)

		req, _ := http.NewRequest("GET", "/admin/racks", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Actions
		controller.ListRacks(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("success with inventory_id", func(t *testing.T) {
		// Set
		mockService := new(MockPalletRackService)
		controller := controllers.NewAdminPalletRackController(mockService)

		inventoryID := uint(1)
		response := &apiContracts.PaginatedRacksResponse{
			Data:  []apiContracts.TransformedRack{},
			Total: 0,
			Page:  1,
			Limit: 10,
		}

		req, _ := http.NewRequest("GET", "/admin/racks?inventory_id=1", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Expectations
		mockService.On("ListRacks", &inventoryID, 1, 10).Return(response, nil)

		// Actions
		controller.ListRacks(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid inventory_id", func(t *testing.T) {
		// Set
		mockService := new(MockPalletRackService)
		controller := controllers.NewAdminPalletRackController(mockService)

		req, _ := http.NewRequest("GET", "/admin/racks?inventory_id=abc", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Actions
		controller.ListRacks(c)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "ListRacks")
	})

	t.Run("service error", func(t *testing.T) {
		// Set
		mockService := new(MockPalletRackService)
		controller := controllers.NewAdminPalletRackController(mockService)

		req, _ := http.NewRequest("GET", "/admin/racks", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Expectations
		mockService.On("ListRacks", (*uint)(nil), 1, 10).Return((*apiContracts.PaginatedRacksResponse)(nil), assert.AnError)

		// Actions
		controller.ListRacks(c)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
