package routes

import (
	baseContainer "go_inventory/Container"
	authControllerPkg "go_inventory/SupplyInventory/Application/Controllers/Auth"
	palletControllerPkg "go_inventory/SupplyInventory/Application/Controllers/Pallet"
	palletRackControllerPkg "go_inventory/SupplyInventory/Application/Controllers/PalletRack"
	palletizedProductControllerPkg "go_inventory/SupplyInventory/Application/Controllers/PalletizedProduct"
	userControllerPkg "go_inventory/SupplyInventory/Application/Controllers/User"
	middlewares "go_inventory/SupplyInventory/Application/Middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, dbInstance *gorm.DB) {
	container := baseContainer.BuildContainer(dbInstance)

	apiV1 := router.Group("/api/v1")

	container.Invoke(func(palletController *palletControllerPkg.PalletController, authMiddleware *middlewares.AuthMiddleware) {
		palletsGroup := apiV1.Group("/pallets")
		palletsGroup.Use(authMiddleware.Handler())
		palletController.Register(palletsGroup)
	})

	container.Invoke(func(ctrl *palletizedProductControllerPkg.PalletizedProductController, authMiddleware *middlewares.AuthMiddleware) {
		palletProductsGroup := apiV1.Group("/pallet/products")
		palletProductsGroup.Use(authMiddleware.Handler())
		ctrl.RegisterProductPallet(palletProductsGroup)
	})

	container.Invoke(func(palletRackController *palletRackControllerPkg.PalletRackController, authMiddleware *middlewares.AuthMiddleware) {
		racksGroup := apiV1.Group("/racks")
		racksGroup.Use(authMiddleware.Handler())
		palletRackController.RegisterPalletRack(racksGroup)
	})

	container.Invoke(func(loginController *authControllerPkg.AuthController) {
		loginController.RegisterLogin(apiV1.Group("/auth"))
	})

	container.Invoke(func(userController *userControllerPkg.UserController) {
		userController.RegisterUserRoutes(apiV1.Group("/users"))
	})
}
