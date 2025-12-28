package container

import (
	authController "go_inventory/SupplyInventory/Application/Controllers/Auth"
	palletController "go_inventory/SupplyInventory/Application/Controllers/Pallet"
	palletRackController "go_inventory/SupplyInventory/Application/Controllers/PalletRack"
	palletizedProductController "go_inventory/SupplyInventory/Application/Controllers/PalletizedProduct"
	userController "go_inventory/SupplyInventory/Application/Controllers/User"
	middlewares "go_inventory/SupplyInventory/Application/Middlewares"

	jwt "go_inventory/SupplyInventory/Application/Services/Jwt"
	pallet "go_inventory/SupplyInventory/Application/Services/Pallet"
	palletrack "go_inventory/SupplyInventory/Application/Services/PalletRack"
	palletizedproduct "go_inventory/SupplyInventory/Application/Services/PalletizedProduct"
	qrcode "go_inventory/SupplyInventory/Application/Services/QrCode"
	user "go_inventory/SupplyInventory/Application/Services/User"

	repositoriesPallet "go_inventory/SupplyInventory/Infrastructure/repositories/Pallet"
	repositoriesPalletRack "go_inventory/SupplyInventory/Infrastructure/repositories/PalletRack"
	repositoriesPalletizedProduct "go_inventory/SupplyInventory/Infrastructure/repositories/PalletizedProduct"
	repositoriesUser "go_inventory/SupplyInventory/Infrastructure/repositories/User"

	contractPallet "go_inventory/SupplyInventory/Domain/contracts/repositories/Pallet"
	contractPalletRack "go_inventory/SupplyInventory/Domain/contracts/repositories/PalletRack"
	contractPalletized "go_inventory/SupplyInventory/Domain/contracts/repositories/PalletizedProduct"
	contractUser "go_inventory/SupplyInventory/Domain/contracts/repositories/User"

	"go.uber.org/dig"
	"gorm.io/gorm"
)

func BuildContainer(db *gorm.DB) *dig.Container {
	container := dig.New()

	// DataBase
	container.Provide(func() *gorm.DB { return db })

	// Repositories
	container.Provide(func(db *gorm.DB) contractPallet.PalletRepository {
		return repositoriesPallet.NewPalletRepository(db)
	})

	container.Provide(func(db *gorm.DB) contractPalletRack.PalletRackRepository {
		return repositoriesPalletRack.NewPalletRackRepository(db)
	})

	container.Provide(func(db *gorm.DB, palletRepo contractPallet.PalletRepository) contractPalletized.PalletizedProductRepository {
		return repositoriesPalletizedProduct.NewPalletizedProductRepository(db, palletRepo)
	})

	container.Provide(func(db *gorm.DB) contractUser.UserRepository {
		return repositoriesUser.NewUserRepository(db)
	})

	// Services
	container.Provide(qrcode.NewQRCodeService)

	container.Provide(func(repo contractPallet.PalletRepository, qrService qrcode.QRCodeService) pallet.PalletService {
		return pallet.NewPalletService(repo, qrService)
	})

	container.Provide(func(repo contractPalletRack.PalletRackRepository) palletrack.PalletRackService {
		return palletrack.NewPalletRackService(repo)
	})

	container.Provide(func(palletRepo contractPallet.PalletRepository, palletizedRepo contractPalletized.PalletizedProductRepository) palletizedproduct.PalletizedProductService {
		return palletizedproduct.NewPalletizedProductService(palletRepo, palletizedRepo)
	})

	container.Provide(jwt.NewJWTService)

	container.Provide(func(repo contractUser.UserRepository) user.UserService {
		return user.NewUserService(repo)
	})

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
