package controllers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	errors "go_inventory/Helpers/Errors"
	pallet "go_inventory/SupplyInventory/Application/Controllers/Pallet"
	entities "go_inventory/SupplyInventory/Domain/Entities"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListPallets_Error(t *testing.T) {

	gin.SetMode(gin.TestMode)
	// Set
	mockService := new(mockPalletService)
	controller := pallet.NewPalletController(mockService)

	// Expectations
	mockService.On("ListPallets", mock.AnythingOfType("*uint"), mock.AnythingOfType("*uint")).Return([]entities.PalletEntity{}, errors.NewAppError("fail", 404))

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	controller.ListPallets(c)

	// Assertions
	assert.Equal(t, 404, w.Code)
}

func TestFindPalletById_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletService)
	controller := pallet.NewPalletController(mockService)
	palletEntity := &entities.PalletEntity{ID: 1, Name: "Pallet1"}
	mockService.On("FindPalletById", uint(1)).Return(palletEntity, (*errors.AppError)(nil))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	controller.FindPalletById(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFindPalletById_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletService)
	controller := pallet.NewPalletController(mockService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "abc"})
	controller.FindPalletById(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFindPalletById_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletService)
	controller := pallet.NewPalletController(mockService)
	mockService.On("FindPalletById", uint(2)).Return(nil, errors.NewAppError("not found", 404))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "2"})
	controller.FindPalletById(c)
	assert.Equal(t, 404, w.Code)
}

func TestDeletePalletById_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletService)
	controller := pallet.NewPalletController(mockService)
	mockService.On("DeletePalletById", uint(1)).Return(true, (*errors.AppError)(nil))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	controller.DeletePalletById(c)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeletePalletById_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletService)
	controller := pallet.NewPalletController(mockService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "abc"})
	controller.DeletePalletById(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeletePalletById_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletService)
	controller := pallet.NewPalletController(mockService)
	mockService.On("DeletePalletById", uint(2)).Return(false, errors.NewAppError("fail", 404))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "2"})
	controller.DeletePalletById(c)
	assert.Equal(t, 404, w.Code)
}

type mockPalletService struct {
	mock.Mock
}

func (m *mockPalletService) ListPallets(palletRackId *uint, productId *uint) ([]entities.PalletEntity, *errors.AppError) {
	args := m.Called(palletRackId, productId)
	return args.Get(0).([]entities.PalletEntity), args.Get(1).(*errors.AppError)
}
func (m *mockPalletService) FindPalletById(id uint) (*entities.PalletEntity, *errors.AppError) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*entities.PalletEntity), args.Get(1).(*errors.AppError)
}
func (m *mockPalletService) DeletePalletById(id uint) (bool, *errors.AppError) {
	args := m.Called(id)
	return args.Bool(0), args.Get(1).(*errors.AppError)
}
func (m *mockPalletService) CreatePallet(name string, rackID uint) (*entities.PalletEntity, *errors.AppError) {
	args := m.Called(name, rackID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*entities.PalletEntity), args.Get(1).(*errors.AppError)
}
func (m *mockPalletService) UpdatePallet(id uint, name string, rackID uint) (*entities.PalletEntity, *errors.AppError) {
	args := m.Called(id, name, rackID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*entities.PalletEntity), args.Get(1).(*errors.AppError)
}
func (m *mockPalletService) GeneratePalletsCsvFile(palletRackId *uint, productId *uint) (string, *errors.AppError) {
	args := m.Called(palletRackId, productId)
	return args.String(0), args.Get(1).(*errors.AppError)
}

func TestListPallets_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set
	mockService := new(mockPalletService)
	controller := pallet.NewPalletController(mockService)
	pallets := []entities.PalletEntity{{ID: 1, Name: "Pallet1"}}

	// Expectations
	mockService.On("ListPallets", mock.AnythingOfType("*uint"), mock.AnythingOfType("*uint")).Return(pallets, (*errors.AppError)(nil))

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	controller.ListPallets(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreatePallet_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletService)
	controller := pallet.NewPalletController(mockService)
	newPallet := &entities.PalletEntity{ID: 1, Name: "Pallet1"}
	mockService.On("CreatePallet", "Pallet1", uint(2)).Return(newPallet, (*errors.AppError)(nil))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/pallets", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(bytes.NewReader([]byte(`{"Name":"Pallet1","PalletRackID":2}`)))
	controller.CreatePallet(c)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreatePallet_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletService)
	controller := pallet.NewPalletController(mockService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/pallets", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(bytes.NewReader([]byte(`{"Name":123}`)))
	controller.CreatePallet(c)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestCreatePallet_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletService)
	controller := pallet.NewPalletController(mockService)
	mockService.On("CreatePallet", "Pallet1", uint(2)).Return(nil, errors.NewAppError("fail", 422))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/pallets", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(bytes.NewReader([]byte(`{"Name":"Pallet1","PalletRackID":2}`)))
	controller.CreatePallet(c)
	assert.Equal(t, 422, w.Code)
}

func TestUpdatePallet_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletService)
	controller := pallet.NewPalletController(mockService)
	updatedPallet := &entities.PalletEntity{ID: 1, Name: "Pallet1"}
	mockService.On("UpdatePallet", uint(1), "Pallet1", uint(2)).Return(updatedPallet, (*errors.AppError)(nil))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	c.Request = httptest.NewRequest("PATCH", "/pallets/1", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(bytes.NewReader([]byte(`{"Name":"Pallet1","PalletRackID":2}`)))
	controller.UpdatePallet(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdatePallet_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletService)
	controller := pallet.NewPalletController(mockService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "abc"})
	c.Request = httptest.NewRequest("PATCH", "/pallets/abc", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(bytes.NewReader([]byte(`{"Name":"Pallet1","PalletRackID":2}`)))
	controller.UpdatePallet(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdatePallet_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletService)
	controller := pallet.NewPalletController(mockService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	c.Request = httptest.NewRequest("PATCH", "/pallets/1", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(bytes.NewReader([]byte(`{"Name":123}`)))
	controller.UpdatePallet(c)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestUpdatePallet_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mockPalletService)
	controller := pallet.NewPalletController(mockService)
	mockService.On("UpdatePallet", uint(1), "Pallet1", uint(2)).Return(nil, errors.NewAppError("fail", 422))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	c.Request = httptest.NewRequest("PATCH", "/pallets/1", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(bytes.NewReader([]byte(`{"Name":"Pallet1","PalletRackID":2}`)))
	controller.UpdatePallet(c)
	assert.Equal(t, 422, w.Code)
}

func TestExportPalletsCsv_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set
	mockService := new(mockPalletService)
	controller := pallet.NewPalletController(mockService)
	url := "http://localhost:3000/reports/Pallets/2023-01-01_12-00-00_123456789.csv"

	// Expectations
	mockService.On("GeneratePalletsCsvFile", mock.AnythingOfType("*uint"), mock.AnythingOfType("*uint")).Return(url, (*errors.AppError)(nil))

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	controller.ExportPalletsCsv(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, url, data["url"])
}

func TestExportPalletsCsv_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set
	mockService := new(mockPalletService)
	controller := pallet.NewPalletController(mockService)

	// Expectations
	mockService.On("GeneratePalletsCsvFile", mock.AnythingOfType("*uint"), mock.AnythingOfType("*uint")).Return("", errors.NewAppError("fail", 404))

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	controller.ExportPalletsCsv(c)

	// Assertions
	assert.Equal(t, 404, w.Code)
}
