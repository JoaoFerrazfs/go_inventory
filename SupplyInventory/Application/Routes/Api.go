package routes

import (
	baseContainer "go_inventory/Container"
	controllers "go_inventory/SupplyInventory/Application/Controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, dbInstance *gorm.DB) {
	container := baseContainer.BuildContainer(dbInstance)

	apiV1 := router.Group("/api/v1")

	container.Invoke(func(palletController *controllers.PalletController) {
		palletController.Register(apiV1.Group("/pallets"))
	})

	err := container.Invoke(func(ctrl *controllers.PalletizedProductController) {
		ctrl.RegisterProductPallet(apiV1.Group("/pallet/products"))
	})
	if err != nil {
		panic(err)
	}

	container.Invoke(func(palletRackController *controllers.PalletRackController) {
		palletRackController.RegisterPalletRack(apiV1.Group("/racks"))
	})
}
