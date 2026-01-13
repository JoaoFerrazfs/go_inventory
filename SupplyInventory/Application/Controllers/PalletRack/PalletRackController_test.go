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
	apiContracts "go_inventory/SupplyInventory/Application/ApiContracts"
	palletRack "go_inventory/SupplyInventory/Application/Controllers/PalletRack"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	testutils "go_inventory/SupplyInventory/tests/testutils"
)

type mockPalletRackService struct {
	mock.Mock
}

func (m *mockPalletRackService) Create(name, location string, totalCapacity int, inventoryID uint) (*entities.PalletRackEntity, error) {
	args := m.Called(name, location, totalCapacity, inventoryID)
	return args.Get(0).(*entities.PalletRackEntity), nil
}
func (m *mockPalletRackService) ListRacks(inventoryID *uint, page int, limit int) (*apiContracts.PaginatedRacksResponse, error) {
	args := m.Called(inventoryID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*apiContracts.PaginatedRacksResponse), args.Error(1)
}
func (m *mockPalletRackService) FindPalletById(id uint) (*entities.PalletRackEntity, *errors.AppError) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*entities.PalletRackEntity), nil
}
func (m *mockPalletRackService) DeleteRack(id uint) (bool, *errors.AppError) {
	args := m.Called(id)
	return args.Bool(0), args.Get(1).(*errors.AppError)
}

func TestCreatePalletRack_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set
	mockService := new(mockPalletRackService)
	controller := palletRack.NewPalletRackController(mockService)
	rack := &entities.PalletRackEntity{ID: 1, Name: "Rack1"}

	// Expectations
	mockService.On("Create", "Rack1", "A1", 100, testutils.DefaultInventoryID).Return(rack, nil)

	// Actions
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("inventoryID", testutils.DefaultInventoryID)
		c.Next()
	})
	group := r.Group("/")
	controller.RegisterPalletRack(group)
	body, _ := json.Marshal(map[string]interface{}{"Name": "Rack1", "Location": "A1", "TotalCapacity": 100})
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreatePalletRack_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletRackService)
	controller := palletRack.NewPalletRackController(mockService)
	r := gin.Default()
	group := r.Group("/")
	controller.RegisterPalletRack(group)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(`{"Name":123}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 422, w.Code)
}

func TestListRacks_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletRackService)
	controller := palletRack.NewPalletRackController(mockService)
	inventoryID := testutils.DefaultInventoryID
	response := &apiContracts.PaginatedRacksResponse{
		Data:  []apiContracts.TransformedRack{{ID: 1, Name: "Rack1"}},
		Total: 1,
		Page:  1,
		Limit: 10,
	}
	mockService.On("ListRacks", &inventoryID, 1, 10).Return(response, nil)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("inventoryID", testutils.DefaultInventoryID)
		c.Next()
	})
	group := r.Group("/")
	controller.RegisterPalletRack(group)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListRacks_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletRackService)
	controller := palletRack.NewPalletRackController(mockService)
	inventoryID := testutils.DefaultInventoryID
	mockService.On("ListRacks", &inventoryID, 1, 10).Return((*apiContracts.PaginatedRacksResponse)(nil), assert.AnError)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("inventoryID", testutils.DefaultInventoryID)
		c.Next()
	})
	group := r.Group("/")
	controller.RegisterPalletRack(group)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestFindRackById_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletRackService)
	controller := palletRack.NewPalletRackController(mockService)
	rack := &entities.PalletRackEntity{ID: 1, Name: "Rack1"}
	mockService.On("FindPalletById", uint(1)).Return(rack, (*errors.AppError)(nil))
	r := gin.Default()
	group := r.Group("/")
	controller.RegisterPalletRack(group)
	req, _ := http.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFindRackById_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletRackService)
	controller := palletRack.NewPalletRackController(mockService)
	mockService.On("FindPalletById", uint(2)).Return(nil, errors.NewAppError("not found", 404))
	r := gin.Default()
	group := r.Group("/")
	controller.RegisterPalletRack(group)
	req, _ := http.NewRequest("GET", "/2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestDeleteRack_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletRackService)
	controller := palletRack.NewPalletRackController(mockService)
	mockService.On("DeleteRack", uint(1)).Return(true, (*errors.AppError)(nil))
	r := gin.Default()
	group := r.Group("/")
	controller.RegisterPalletRack(group)
	req, _ := http.NewRequest("DELETE", "/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteRack_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletRackService)
	controller := palletRack.NewPalletRackController(mockService)
	mockService.On("DeleteRack", uint(2)).Return(false, errors.NewAppError("fail", 404))
	r := gin.Default()
	group := r.Group("/")
	controller.RegisterPalletRack(group)
	req, _ := http.NewRequest("DELETE", "/2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}
