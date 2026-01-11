package controllers

import (
	"fmt"
	"net/http"

	errors "go_inventory/Helpers/Errors"
	requestsHelper "go_inventory/Helpers/RequestsHelper"
	inventoryRequest "go_inventory/SupplyInventory/Application/Requests/Inventory"
	inventoryService "go_inventory/SupplyInventory/Application/Services/Inventory"
	userService "go_inventory/SupplyInventory/Application/Services/User"

	"github.com/gin-gonic/gin"
)

type InventoryController struct {
	inventoryService inventoryService.InventoryService
	userService      userService.UserService
}

func NewInventoryController(inventoryService inventoryService.InventoryService, userService userService.UserService) *InventoryController {
	return &InventoryController{
		inventoryService: inventoryService,
		userService:      userService,
	}
}

func (controller *InventoryController) Register(group *gin.RouterGroup) {
	group.GET("/", controller.ListInventories)
	group.GET("/:id", controller.GetInventoryById)
	group.POST("/", controller.CreateInventory)
	group.PUT("/:id", controller.UpdateInventory)
}

// @Summary List all inventories
// @Tags Inventory
// @Accept json
// @Produce json
// @Success 200 {object} apiContracts.InventoryListResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/inventories [get]
func (controller *InventoryController) ListInventories(c *gin.Context) {
	inventories, appError := controller.inventoryService.ListInventories()
	if appError != nil {
		c.JSON(appError.Code, gin.H{"error": appError.Message})
		return
	}
	c.JSON(200, gin.H{"data": gin.H{"inventories": inventories}})
}

// @Summary Get inventory by ID
// @Tags Inventory
// @Accept json
// @Produce json
// @Param id path int true "Inventory ID"
// @Success 200 {object} apiContracts.InventoryResponse
// @Failure 422 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/inventories/{id} [get]
func (controller *InventoryController) GetInventoryById(c *gin.Context) {
	idParam := c.Param("id")

	var id uint
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		c.JSON(422, gin.H{"error": "Invalid inventory ID"})
		return
	}

	inventory, appError := controller.inventoryService.GetInventoryByID(id)
	if appError != nil {
		c.JSON(appError.Code, gin.H{"error": appError.Message})
		return
	}
	c.JSON(200, gin.H{"data": gin.H{"inventory": inventory}})
}

// @Summary Create a new inventory
// @Tags Inventory
// @Accept json
// @Produce json
// @Param inventory body inventory.InventoryRequest true "Inventory Data"
// @Success 201 {object} apiContracts.InventoryResponse
// @Failure 422 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/inventories [post]
func (controller *InventoryController) CreateInventory(c *gin.Context) {
	var req inventoryRequest.InventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(422, requestsHelper.FormatValidationErrors(err))
		return
	}

	userIDInterface, exists := c.Get("userID")
	if !exists {
		appError := errors.NewAppError("User not authenticated", 401)
		c.JSON(appError.Code, requestsHelper.FormatValidationErrors(appError))
		return
	}

	userID := userIDInterface.(uint)

	user, appErr := controller.userService.GetUserByID(userID)
	if appErr != nil {
		c.JSON(appErr.Code, requestsHelper.FormatValidationErrors(appErr))
		return
	}

	inventory, appError := controller.inventoryService.CreateInventory(req.Name, req.Description, *user)
	if appError != nil {
		c.JSON(appError.Code, requestsHelper.FormatValidationErrors(appError))
		return
	}

	c.JSON(201, gin.H{"data": gin.H{"inventory": inventory}})
}

func (controller *InventoryController) UpdateInventory(c *gin.Context) {
	var req inventoryRequest.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(422, requestsHelper.FormatValidationErrors(err))
		return
	}

	inventoryIdParam := c.Param("id")
	var inventoryId uint
	_, err := fmt.Sscanf(inventoryIdParam, "%d", &inventoryId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid inventory ID"})
		return
	}

	if req.Name == "" && req.Description == "" && req.Status == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "At least one field (name, description, status) must be provided"})
		return
	}

	inventory, appError := controller.inventoryService.UpdateInventory(inventoryId, req.Name, req.Description, req.Status)
	if appError != nil {
		c.JSON(appError.Code, gin.H{"error": appError.Message})
		return
	}

	c.JSON(200, gin.H{"data": gin.H{"inventory": inventory}})
}
