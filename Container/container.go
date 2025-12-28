package container

import (
	authController "go_inventory/SupplyInventory/Application/Controllers/Auth"
	palletController "go_inventory/SupplyInventory/Application/Controllers/Pallet"
	palletRackController "go_inventory/SupplyInventory/Application/Controllers/PalletRack"
	palletizedProductController "go_inventory/SupplyInventory/Application/Controllers/PalletizedProduct"
	userController "go_inventory/SupplyInventory/Application/Controllers/User"
	middlewares "go_inventory/SupplyInventory/Application/Middlewares"
	services "go_inventory/SupplyInventory/Application/Services"

	repositoriesPallet "go_inventory/SupplyInventory/Infrastructure/repositories/Pallet"
	repositoriesPalletRack "go_inventory/SupplyInventory/Infrastructure/repositories/PalletRack"
	repositoriesPalletizedProduct "go_inventory/SupplyInventory/Infrastructure/repositories/PalletizedProduct"
	repositoriesUser "go_inventory/SupplyInventory/Infrastructure/repositories/User"

	"go.uber.org/dig"
	"gorm.io/gorm"
)

func BuildContainer(db *gorm.DB) *dig.Container {
	container := dig.New()

	// DataBase
	container.Provide(func() *gorm.DB { return db })

	// Repositories
	container.Provide(repositoriesPallet.NewPalletRepository)
	container.Provide(repositoriesPalletRack.NewPalletRackRepository)
	container.Provide(repositoriesPalletizedProduct.NewPalletizedProductRepository)
	container.Provide(repositoriesUser.NewUserRepository)

	// Services
	container.Provide(services.NewQRCodeService)
	container.Provide(services.NewPalletService)
	container.Provide(services.NewPalletRackService)
	container.Provide(services.NewPalletizedProductService)
	container.Provide(services.NewJWTService)
	container.Provide(services.NewUserService)

	// Api Controllers
	container.Provide(palletController.NewPalletController)
	container.Provide(palletizedProductController.NewPalletizedProductController)
	container.Provide(palletRackController.NewPalletRackController)
	container.Provide(authController.NewAuthController)
	container.Provide(userController.NewUserController)

	// Middlewares
	container.Provide(middlewares.NewAuthMiddleware)

	return container
}
