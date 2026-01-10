package inventory_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	appErrors "go_inventory/Helpers/Errors"
	inventory "go_inventory/SupplyInventory/Application/Controllers/Inventory"
	inventoryRequest "go_inventory/SupplyInventory/Application/Requests/Inventory"
	inventoryService "go_inventory/SupplyInventory/Application/Services/Inventory"
	userService "go_inventory/SupplyInventory/Application/Services/User"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	"go_inventory/SupplyInventory/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListInventories_Success(t *testing.T) {
	// Set
	inventoryRepo := new(mocks.InventoryRepository)
	userRepo := new(mocks.UserRepository)
	inventorySvc := inventoryService.NewInventoryService(inventoryRepo)
	userSvc := userService.NewUserService(userRepo)
	controller := inventory.NewInventoryController(inventorySvc, userSvc)

	inventories := []entities.InventoryEntity{
		{ID: 1, Name: "Inventory 1"},
		{ID: 2, Name: "Inventory 2"},
	}

	// Expectations
	inventoryRepo.On("Where", mock.Anything).Return(inventories, (*appErrors.AppError)(nil))

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	controller.ListInventories(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response["data"])
	inventoryRepo.AssertExpectations(t)
}

func TestListInventories_Error(t *testing.T) {
	// Set
	inventoryRepo := new(mocks.InventoryRepository)
	userRepo := new(mocks.UserRepository)
	inventorySvc := inventoryService.NewInventoryService(inventoryRepo)
	userSvc := userService.NewUserService(userRepo)
	controller := inventory.NewInventoryController(inventorySvc, userSvc)

	// Expectations
	inventoryRepo.On("Where", mock.Anything).Return(nil, appErrors.NewAppError("Database error", 500))

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	controller.ListInventories(c)

	// Assertions
	assert.Equal(t, 500, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Database error", response["error"])
	inventoryRepo.AssertExpectations(t)
}

func TestGetInventoryById_Success(t *testing.T) {
	// Set
	inventoryRepo := new(mocks.InventoryRepository)
	userRepo := new(mocks.UserRepository)
	inventorySvc := inventoryService.NewInventoryService(inventoryRepo)
	userSvc := userService.NewUserService(userRepo)
	controller := inventory.NewInventoryController(inventorySvc, userSvc)

	// Expectations
	inventoryRepo.On("FindById", uint(1)).Return(&entities.InventoryEntity{ID: 1, Name: "Inventory 1"}, (*appErrors.AppError)(nil))

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
	inventoryRepo.AssertExpectations(t)
}

func TestGetInventoryById_InvalidID(t *testing.T) {
	// Set
	inventoryRepo := new(mocks.InventoryRepository)
	userRepo := new(mocks.UserRepository)
	inventorySvc := inventoryService.NewInventoryService(inventoryRepo)
	userSvc := userService.NewUserService(userRepo)
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
	inventoryRepo := new(mocks.InventoryRepository)
	userRepo := new(mocks.UserRepository)
	inventorySvc := inventoryService.NewInventoryService(inventoryRepo)
	userSvc := userService.NewUserService(userRepo)
	controller := inventory.NewInventoryController(inventorySvc, userSvc)

	// Expectations
	inventoryRepo.On("FindById", uint(1)).Return(nil, appErrors.NewAppError("Inventory not found", 404))

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
	inventoryRepo.AssertExpectations(t)
}

func TestCreateInventory_Success(t *testing.T) {
	// Set
	inventoryRepo := new(mocks.InventoryRepository)
	userRepo := new(mocks.UserRepository)
	inventorySvc := inventoryService.NewInventoryService(inventoryRepo)
	userSvc := userService.NewUserService(userRepo)
	controller := inventory.NewInventoryController(inventorySvc, userSvc)

	req := inventoryRequest.InventoryRequest{Name: "New Inventory", Description: "Description"}
	user := &entities.UserEntity{ID: 1, Name: "User"}

	// Expectations
	userRepo.On("FindByID", uint(1)).Return(user, nil)
	inventoryRepo.On("Create", mock.AnythingOfType("*entities.InventoryEntity")).Return((*appErrors.AppError)(nil))

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
	inventoryRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestCreateInventory_InvalidJSON(t *testing.T) {
	// Set
	inventoryRepo := new(mocks.InventoryRepository)
	userRepo := new(mocks.UserRepository)
	inventorySvc := inventoryService.NewInventoryService(inventoryRepo)
	userSvc := userService.NewUserService(userRepo)
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
	inventoryRepo := new(mocks.InventoryRepository)
	userRepo := new(mocks.UserRepository)
	inventorySvc := inventoryService.NewInventoryService(inventoryRepo)
	userSvc := userService.NewUserService(userRepo)
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
	inventoryRepo := new(mocks.InventoryRepository)
	userRepo := new(mocks.UserRepository)
	inventorySvc := inventoryService.NewInventoryService(inventoryRepo)
	userSvc := userService.NewUserService(userRepo)
	controller := inventory.NewInventoryController(inventorySvc, userSvc)

	req := inventoryRequest.InventoryRequest{Name: "New Inventory"}

	// Expectations
	userRepo.On("FindByID", uint(1)).Return(nil, errors.New("user not found"))

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
	userRepo.AssertExpectations(t)
}

func TestCreateInventory_ErrorCreating(t *testing.T) {
	// Set
	inventoryRepo := new(mocks.InventoryRepository)
	userRepo := new(mocks.UserRepository)
	inventorySvc := inventoryService.NewInventoryService(inventoryRepo)
	userSvc := userService.NewUserService(userRepo)
	controller := inventory.NewInventoryController(inventorySvc, userSvc)

	req := inventoryRequest.InventoryRequest{Name: "New Inventory"}
	user := &entities.UserEntity{ID: 1, Name: "User"}

	// Expectations
	userRepo.On("FindByID", uint(1)).Return(user, nil)
	inventoryRepo.On("Create", mock.AnythingOfType("*entities.InventoryEntity")).Return(appErrors.NewAppError("Creation error", 500))

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
	inventoryRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}
