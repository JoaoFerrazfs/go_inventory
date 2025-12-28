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

	"go.uber.org/fx"
	"gorm.io/gorm"
)

func BuildOptions(db *gorm.DB) fx.Option {
	return fx.Options(
		// Database
		fx.Supply(db),

		// Repositories Module
		fx.Module("repositories",
			fx.Provide(func(db *gorm.DB) contractPallet.PalletRepository {
				return repositoriesPallet.NewPalletRepository(db)
			}),
			fx.Provide(func(db *gorm.DB) contractPalletRack.PalletRackRepository {
				return repositoriesPalletRack.NewPalletRackRepository(db)
			}),
			fx.Provide(func(db *gorm.DB, palletRepo contractPallet.PalletRepository) contractPalletized.PalletizedProductRepository {
				return repositoriesPalletizedProduct.NewPalletizedProductRepository(db, palletRepo)
			}),
			fx.Provide(func(db *gorm.DB) contractUser.UserRepository {
				return repositoriesUser.NewUserRepository(db)
			}),
		),

		// Services Module
		fx.Module("services",
			fx.Provide(qrcode.NewQRCodeService),
			fx.Provide(func(repo contractPallet.PalletRepository, qrService qrcode.QRCodeService) pallet.PalletService {
				return pallet.NewPalletService(repo, qrService)
			}),
			fx.Provide(func(repo contractPalletRack.PalletRackRepository) palletrack.PalletRackService {
				return palletrack.NewPalletRackService(repo)
			}),
			fx.Provide(func(palletRepo contractPallet.PalletRepository, palletizedRepo contractPalletized.PalletizedProductRepository) palletizedproduct.PalletizedProductService {
				return palletizedproduct.NewPalletizedProductService(palletRepo, palletizedRepo)
			}),
			fx.Provide(jwt.NewJWTService),
			fx.Provide(func(repo contractUser.UserRepository) user.UserService {
				return user.NewUserService(repo)
			}),
		),

		// Controllers Module
		fx.Module("controllers",
			fx.Provide(palletController.NewPalletController),
			fx.Provide(palletizedProductController.NewPalletizedProductController),
			fx.Provide(palletRackController.NewPalletRackController),
			fx.Provide(authController.NewAuthController),
			fx.Provide(userController.NewUserController),
		),

		// Middleware Module
		fx.Module("middleware",
			fx.Provide(middlewares.NewAuthMiddleware),
		),
	)
}
