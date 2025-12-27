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
	domain "go_inventory/SupplyInventory/Domain"
)

type mockPalletRackService struct {
	mock.Mock
}

func (m *mockPalletRackService) Create(name, location string, totalCapacity int) (*domain.PalletRackEntity, error) {
	args := m.Called(name, location, totalCapacity)
	return args.Get(0).(*domain.PalletRackEntity), nil
}
func (m *mockPalletRackService) ListRacks() ([]apiContracts.TransformedRack, error) {
	args := m.Called()
	return args.Get(0).([]apiContracts.TransformedRack), args.Error(1)
}
func (m *mockPalletRackService) FindPalletById(id uint) (*domain.PalletRackEntity, *errors.AppError) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*domain.PalletRackEntity), nil
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
	rack := &domain.PalletRackEntity{ID: 1, Name: "Rack1"}

	// Expectations
	mockService.On("Create", "Rack1", "A1", 100).Return(rack, nil)

	// Actions
	r := gin.Default()
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
	racks := []apiContracts.TransformedRack{{ID: 1, Name: "Rack1"}}
	mockService.On("ListRacks").Return(racks, nil)
	r := gin.Default()
	group := r.Group("/")
	controller.RegisterPalletRack(group)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestListRacks_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletRackService)
	controller := palletRack.NewPalletRackController(mockService)
	mockService.On("ListRacks").Return([]apiContracts.TransformedRack{}, assert.AnError)
	r := gin.Default()
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
	rack := &domain.PalletRackEntity{ID: 1, Name: "Rack1"}
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
