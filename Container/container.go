package container

import (
	controllers "go_inventory/SupplyInventory/Application/Controllers"
	services "go_inventory/SupplyInventory/Application/Services"
	infrastructure "go_inventory/SupplyInventory/Infrastructure"

	"go.uber.org/dig"
	"gorm.io/gorm"
)

func BuildContainer(db *gorm.DB) *dig.Container {
	c := dig.New()

	// Registrar dependências
	c.Provide(func() *gorm.DB { return db })

	// Repositories
	c.Provide(infrastructure.NewPalletRepository)
	c.Provide(infrastructure.NewPalletRackRepository)

	// Services
	c.Provide(services.NewQRCodeService)
	c.Provide(services.NewPalletService)
	c.Provide(services.NewPalletRackService)

	// Controllers
	c.Provide(controllers.NewPalletController)
	c.Provide(controllers.NewPalletizedProductController)
	c.Provide(controllers.NewPalletRackController)

	return c
}
