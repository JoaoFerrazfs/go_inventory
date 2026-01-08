package middlewares

import (
	"strconv"

	errors "go_inventory/Helpers/Errors"
	inventoryRepository "go_inventory/SupplyInventory/Domain/contracts/repositories/Inventory"

	"github.com/gin-gonic/gin"
)

type InventoryMiddleware struct {
	InventoryRepository inventoryRepository.InventoryRepository
}

func NewInventoryMiddleware(repo inventoryRepository.InventoryRepository) *InventoryMiddleware {
	return &InventoryMiddleware{InventoryRepository: repo}
}

func (m *InventoryMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("X-Inventory-ID")
		if header == "" {
			c.AbortWithStatusJSON(422, errors.NewAppError("X-Inventory-ID header is required", 422))
			return
		}

		id64, err := strconv.ParseUint(header, 10, 32)
		if err != nil {
			c.AbortWithStatusJSON(422, errors.NewAppError("X-Inventory-ID must be a valid number", 422))
			return
		}

		id := uint(id64)
		exists, appErr := m.InventoryRepository.Exists(id)
		if appErr != nil || !exists {
			c.AbortWithStatusJSON(appErr.Code, appErr)
			return
		}

		c.Set("inventoryID", id)
		c.Next()
	}
}

func GetInventoryID(c *gin.Context) uint {
	if v, ok := c.Get("inventoryID"); ok {
		if id, ok2 := v.(uint); ok2 {
			return id
		}
	}
	return 0
}
