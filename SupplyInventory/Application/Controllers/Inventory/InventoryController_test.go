package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	appErrors "go_inventory/Helpers/Errors"
	inventory "go_inventory/SupplyInventory/Application/Controllers/Inventory"
	inventoryRequest "go_inventory/SupplyInventory/Application/Requests/Inventory"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	"go_inventory/SupplyInventory/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestListInventories_Success(t *testing.T) {
	// Set
	inventorySvc := new(mocks.InventoryService)
	userSvc := new(mocks.UserService)
	controller := inventory.NewInventoryController(inventorySvc, userSvc)

	inventories := []entities.InventoryEntity{
		{ID: 1, Name: "Inventory 1"},
		{ID: 2, Name: "Inventory 2"},
	}

	// Expectations
	inventorySvc.On("ListInventories").Return(inventories, (*appErrors.AppError)(nil))

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	controller.ListInventories(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response["data"])
	inventorySvc.AssertExpectations(t)
}

func TestListInventories_Error(t *testing.T) {
	// Set
	inventorySvc := new(mocks.InventoryService)
	userSvc := new(mocks.UserService)
	controller := inventory.NewInventoryController(inventorySvc, userSvc)

	// Expectations
	inventorySvc.On("ListInventories").Return(nil, appErrors.NewAppError("Database error", 500))

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	controller.ListInventories(c)

	// Assertions
	assert.Equal(t, 500, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Database error", response["error"])
	inventorySvc.AssertExpectations(t)
}

func TestGetInventoryById_Success(t *testing.T) {
	// Set
	inventorySvc := new(mocks.InventoryService)
	userSvc := new(mocks.UserService)
	controller := inventory.NewInventoryController(inventorySvc, userSvc)

	// Expectations
	inventorySvc.On("GetInventoryByID", uint(1)).Return(entities.InventoryEntity{ID: 1, Name: "Inventory 1"}, (*appErrors.AppError)(nil))

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	controller.GetInventoryById(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response["data"])
	inventorySvc.AssertExpectations(t)
}

func TestGetInventoryById_InvalidID(t *testing.T) {
	// Set
	inventorySvc := new(mocks.InventoryService)
	userSvc := new(mocks.UserService)
	controller := inventory.NewInventoryController(inventorySvc, userSvc)

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}
	controller.GetInventoryById(c)

	// Assertions
	assert.Equal(t, 422, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid inventory ID", response["error"])
}

func TestGetInventoryById_NotFound(t *testing.T) {
	// Set
	inventorySvc := new(mocks.InventoryService)
	userSvc := new(mocks.UserService)
	controller := inventory.NewInventoryController(inventorySvc, userSvc)

	// Expectations
	inventorySvc.On("GetInventoryByID", uint(1)).Return(entities.InventoryEntity{}, appErrors.NewAppError("Inventory not found", 404))

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	controller.GetInventoryById(c)

	// Assertions
	assert.Equal(t, 404, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Inventory not found", response["error"])
	inventorySvc.AssertExpectations(t)
}

func TestCreateInventory_Success(t *testing.T) {
	// Set
	inventorySvc := new(mocks.InventoryService)
	userSvc := new(mocks.UserService)
	controller := inventory.NewInventoryController(inventorySvc, userSvc)

	req := inventoryRequest.InventoryRequest{Name: "New Inventory", Description: "Description"}
	user := &entities.UserEntity{ID: 1, Name: "User"}

	// Expectations
	userSvc.On("GetUserByID", uint(1)).Return(user, (*appErrors.AppError)(nil))
	inventorySvc.On("CreateInventory", "New Inventory", "Description", *user).Return(entities.InventoryEntity{ID: 1, Name: "New Inventory", Description: "Description", User: *user}, (*appErrors.AppError)(nil))

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", uint(1))
	body, _ := json.Marshal(req)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	controller.CreateInventory(c)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response["data"])
	inventorySvc.AssertExpectations(t)
	userSvc.AssertExpectations(t)
}

func TestCreateInventory_InvalidJSON(t *testing.T) {
	// Set
	inventorySvc := new(mocks.InventoryService)
	userSvc := new(mocks.UserService)
	controller := inventory.NewInventoryController(inventorySvc, userSvc)

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", uint(1))
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")
	controller.CreateInventory(c)

	// Assertions
	assert.Equal(t, 422, w.Code)
}

func TestCreateInventory_UserNotAuthenticated(t *testing.T) {
	// Set
	inventorySvc := new(mocks.InventoryService)
	userSvc := new(mocks.UserService)
	controller := inventory.NewInventoryController(inventorySvc, userSvc)

	req := inventoryRequest.InventoryRequest{Name: "New Inventory"}

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body, _ := json.Marshal(req)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	controller.CreateInventory(c)

	// Assertions
	assert.Equal(t, 401, w.Code)
}

func TestCreateInventory_UserNotFound(t *testing.T) {
	// Set
	inventorySvc := new(mocks.InventoryService)
	userSvc := new(mocks.UserService)
	controller := inventory.NewInventoryController(inventorySvc, userSvc)

	req := inventoryRequest.InventoryRequest{Name: "New Inventory"}

	// Expectations
	userSvc.On("GetUserByID", uint(1)).Return(nil, appErrors.NewAppError("User not found", 404))

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", uint(1))
	body, _ := json.Marshal(req)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	controller.CreateInventory(c)

	// Assertions
	assert.Equal(t, 404, w.Code)
	userSvc.AssertExpectations(t)
}

func TestCreateInventory_ErrorCreating(t *testing.T) {
	// Set
	inventorySvc := new(mocks.InventoryService)
	userSvc := new(mocks.UserService)
	controller := inventory.NewInventoryController(inventorySvc, userSvc)

	req := inventoryRequest.InventoryRequest{Name: "New Inventory"}
	user := &entities.UserEntity{ID: 1, Name: "User"}

	// Expectations
	userSvc.On("GetUserByID", uint(1)).Return(user, (*appErrors.AppError)(nil))
	inventorySvc.On("CreateInventory", "New Inventory", "", *user).Return(entities.InventoryEntity{}, appErrors.NewAppError("Creation error", 500))

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", uint(1))
	body, _ := json.Marshal(req)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	controller.CreateInventory(c)

	// Assertions
	assert.Equal(t, 500, w.Code)
	inventorySvc.AssertExpectations(t)
	userSvc.AssertExpectations(t)
}
