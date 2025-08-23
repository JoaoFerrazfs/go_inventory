package routes

import (
	controllers "go_inventory/SupplyInventory/Application/Controllers"
	services "go_inventory/SupplyInventory/Application/Services"
	infrastructure "go_inventory/SupplyInventory/Infrastructure"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, dbInstance *gorm.DB) {
	// Criando repositório
	palletRepo := infrastructure.NewPalletRepository(dbInstance)

	// Criando serviço de QR code
	qrService := services.NewQRCodeService()

	// Criando service
	palletService := services.NewPalletService(palletRepo, qrService)

	// Criando controllers
	palletController := controllers.NewPalletController(palletService)

	// Grupo de rotas
	apiV1 := router.Group("/api/v1/pallets")

	// Registrando rotas
	palletController.Register(apiV1)              // rotas GET, POST
	palletController.RegisterProductPallet(apiV1) // rota PATCH /:pallet/products
}
