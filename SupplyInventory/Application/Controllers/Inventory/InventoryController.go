package inventory

import (
	"fmt"

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

func (controller *InventoryController) ListInventories(c *gin.Context) {
	inventories, appError := controller.inventoryService.ListInventories()
	if appError != nil {
		c.JSON(appError.Code, gin.H{"error": appError.Message})
		return
	}
	c.JSON(200, gin.H{"data": gin.H{"inventories": inventories}})
}

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

func (controller *InventoryController) CreateInventory(c *gin.Context) {
	var req inventoryRequest.InventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(422, requestsHelper.FormatValidationErrors(err))
		return
	}

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDInterface.(uint)

	user, appErr := controller.userService.GetUserByID(userID)
	if appErr != nil {
		c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		return
	}

	inventory, appError := controller.inventoryService.CreateInventory(req.Name, req.Description, *user)
	if appError != nil {
		c.JSON(appError.Code, gin.H{"error": appError.Message})
		return
	}
	c.JSON(201, gin.H{"data": gin.H{"inventory": inventory}})
}

func (controller *InventoryController) UpdateInventory(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Update inventory"})
}
