package routes

import (
	"github.com/gin-gonic/gin"

	"go_inventory/SupplyInventory/Application/Controllers"
)

func RegisterRoutes(router *gin.Engine) {
	apiV1 := router.Group("/api/v1/positions")
	controllers.Register(apiV1)
}