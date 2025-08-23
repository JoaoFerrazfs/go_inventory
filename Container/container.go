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
	c.Provide(infrastructure.NewPalletRepository)
	c.Provide(services.NewQRCodeService)
	c.Provide(services.NewPalletService)
	c.Provide(controllers.NewPalletController)

	return c
}
