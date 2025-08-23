package routes

import (
	container "go_inventory/Container"
	controllers "go_inventory/SupplyInventory/Application/Controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, dbInstance *gorm.DB) {
	c := container.BuildContainer(dbInstance)

	c.Invoke(func(palletController *controllers.PalletController) {
		apiV1 := router.Group("/api/v1/pallets")
		palletController.Register(apiV1)
		palletController.RegisterProductPallet(apiV1)
	})
}
