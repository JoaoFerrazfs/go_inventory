package routes

import (
	"github.com/gin-gonic/gin"

	controllers "go_inventory/SupplyInventory/Application/Controllers"
)

func RegisterRoutes(router *gin.Engine) {
	apiV1 := router.Group("/api/v1/pallets")

	controllers.Register(apiV1)
	controllers.RegisterProductPallet(apiV1)
}
