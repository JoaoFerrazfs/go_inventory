package inventory

import "github.com/gin-gonic/gin"

type InventoryController struct{}

func NewInventoryController() *InventoryController {
	return &InventoryController{}
}

func (controller *InventoryController) Register(group *gin.RouterGroup) {
	group.GET("/", controller.ListInventories)
	group.GET("/:id", controller.GetInventoryById)
	group.POST("/", controller.CreateInventory)
	group.PUT("/:id", controller.UpdateInventory)
}

func (controller *InventoryController) ListInventories(c *gin.Context) {
	c.JSON(200, gin.H{"message": "List of inventories"})
}

func (controller *InventoryController) GetInventoryById(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get inventory by ID"})
}

func (controller *InventoryController) CreateInventory(c *gin.Context) {
	c.JSON(201, gin.H{"message": "Create new inventory"})
}

func (controller *InventoryController) UpdateInventory(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Update inventory"})
}
