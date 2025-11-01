package container

import (
	controllers "go_inventory/SupplyInventory/Application/Controllers"
	middlewares "go_inventory/SupplyInventory/Application/Middlewares"
	services "go_inventory/SupplyInventory/Application/Services"
	infrastructure "go_inventory/SupplyInventory/Infrastructure"

	"go.uber.org/dig"
	"gorm.io/gorm"
)

func BuildContainer(db *gorm.DB) *dig.Container {
	container := dig.New()

	// DataBase
	container.Provide(func() *gorm.DB { return db })

	// Repositories
	container.Provide(infrastructure.NewPalletRepository)
	container.Provide(infrastructure.NewPalletRackRepository)
	container.Provide(infrastructure.NewPalletizedProductRepository)

	// Services
	container.Provide(services.NewQRCodeService)
	container.Provide(services.NewPalletService)
	container.Provide(services.NewPalletRackService)
	container.Provide(services.NewPalletizedProductService)
	container.Provide(services.NewJWTService)

	// Api Controllers
	container.Provide(controllers.NewPalletController)
	container.Provide(controllers.NewPalletizedProductController)
	container.Provide(controllers.NewPalletRackController)
	container.Provide(controllers.NewLoginController)

	// Middlewares
	container.Provide(middlewares.NewAuthMiddleware)

	return container
}
