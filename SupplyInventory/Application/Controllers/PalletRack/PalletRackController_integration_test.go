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
	palletRack "go_inventory/SupplyInventory/Application/Controllers/PalletRack"
	requests "go_inventory/SupplyInventory/Application/Requests"
	services "go_inventory/SupplyInventory/Application/Services"
	entities "go_inventory/SupplyInventory/Domain/Entities"
)

type mockPalletRackRepository struct {
	mock.Mock
}

func (m *mockPalletRackRepository) Create(name string, location string, totalCapacity int) (*entities.PalletRackEntity, error) {
	args := m.Called(name, location, totalCapacity)
	return args.Get(0).(*entities.PalletRackEntity), args.Error(1)
}
func (m *mockPalletRackRepository) ListRacks() ([]entities.PalletRackEntity, error) {
	args := m.Called()
	return args.Get(0).([]entities.PalletRackEntity), args.Error(1)
}
func (m *mockPalletRackRepository) FindPalletById(id uint) (*entities.PalletRackEntity, *errors.AppError) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*entities.PalletRackEntity), args.Get(1).(*errors.AppError)
}
func (m *mockPalletRackRepository) DeleteRack(id uint) (bool, *errors.AppError) {
	args := m.Called(id)
	return args.Bool(0), args.Get(1).(*errors.AppError)
}

func setupPalletRackRouter() *gin.Engine {
	mockRepo := new(mockPalletRackRepository)
	service := services.NewPalletRackService(mockRepo)
	controller := palletRack.NewPalletRackController(service)

	r := gin.Default()
	api := r.Group("/api/v1/racks")
	controller.RegisterPalletRack(api)
	return r
}

func TestIntegration_CreatePalletRack(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set
	mockRepo := new(mockPalletRackRepository)
	service := services.NewPalletRackService(mockRepo)
	controller := palletRack.NewPalletRackController(service)
	r := gin.Default()
	api := r.Group("/api/v1/racks")
	controller.RegisterPalletRack(api)

	createReq := requests.PalletRackRequest{
		Name:         "Rack1",
		Location:     "A1",
		TotalCapacity: 100,
	}
	body, _ := json.Marshal(createReq)

	// Expectations
	mockRepo.On("Create", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
		Return(&entities.PalletRackEntity{}, nil)

	// Actions
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/racks/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Assertions
	assert.Contains(t, []int{http.StatusCreated, http.StatusUnprocessableEntity}, w.Code)
}
