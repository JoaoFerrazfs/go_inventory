package container

import (
	"os"

	authController "go_inventory/SupplyInventory/Application/Controllers/Auth"
	inventoryController "go_inventory/SupplyInventory/Application/Controllers/Inventory"
	palletController "go_inventory/SupplyInventory/Application/Controllers/Pallet"
	palletRackController "go_inventory/SupplyInventory/Application/Controllers/PalletRack"
	palletizedProductController "go_inventory/SupplyInventory/Application/Controllers/PalletizedProduct"
	productController "go_inventory/SupplyInventory/Application/Controllers/Product"
	userController "go_inventory/SupplyInventory/Application/Controllers/User"
	middlewares "go_inventory/SupplyInventory/Application/Middlewares"

	inventory "go_inventory/SupplyInventory/Application/Services/Inventory"
	jwt "go_inventory/SupplyInventory/Application/Services/Jwt"
	pallet "go_inventory/SupplyInventory/Application/Services/Pallet"
	palletrack "go_inventory/SupplyInventory/Application/Services/PalletRack"
	palletizedproduct "go_inventory/SupplyInventory/Application/Services/PalletizedProduct"
	productService "go_inventory/SupplyInventory/Application/Services/Product"
	qrcode "go_inventory/SupplyInventory/Application/Services/QrCode"
	user "go_inventory/SupplyInventory/Application/Services/User"

	repositoriesPallet "go_inventory/SupplyInventory/Infrastructure/repositories/Pallet"
	repositoriesPalletRack "go_inventory/SupplyInventory/Infrastructure/repositories/PalletRack"
	repositoriesPalletizedProduct "go_inventory/SupplyInventory/Infrastructure/repositories/PalletizedProduct"
	repositoryProduct "go_inventory/SupplyInventory/Infrastructure/repositories/Product"
	repositoriesUser "go_inventory/SupplyInventory/Infrastructure/repositories/User"

	contractInventory "go_inventory/SupplyInventory/Domain/contracts/repositories/Inventory"
	contractPallet "go_inventory/SupplyInventory/Domain/contracts/repositories/Pallet"
	contractPalletRack "go_inventory/SupplyInventory/Domain/contracts/repositories/PalletRack"
	contractPalletized "go_inventory/SupplyInventory/Domain/contracts/repositories/PalletizedProduct"
	contractProduct "go_inventory/SupplyInventory/Domain/contracts/repositories/Product"
	contractUser "go_inventory/SupplyInventory/Domain/contracts/repositories/User"

	storage "go_inventory/SupplyInventory/Application/Services/Storage"
	dbadapter "go_inventory/SupplyInventory/Infrastructure/repositories/db"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"gorm.io/gorm"

	repositoriesInventory "go_inventory/SupplyInventory/Infrastructure/repositories/Inventory"
)

func BuildOptions(db *gorm.DB) fx.Option {
	return fx.Options(

		fx.WithLogger(func() fxevent.Logger { return fxevent.NopLogger }),
		fx.Supply(db),

		// Repositories
		fx.Provide(func(db *gorm.DB) contractPallet.PalletRepository {
			return repositoriesPallet.NewPalletRepository(dbadapter.NewGormAdapter(db))
		}),
		fx.Provide(func(db *gorm.DB) contractPalletRack.PalletRackRepository {
			return repositoriesPalletRack.NewPalletRackRepository(dbadapter.NewGormAdapter(db))
		}),
		fx.Provide(func(db *gorm.DB, palletRepo contractPallet.PalletRepository) contractPalletized.PalletizedProductRepository {
			return repositoriesPalletizedProduct.NewPalletizedProductRepository(dbadapter.NewGormAdapter(db), palletRepo)
		}),
		fx.Provide(func(db *gorm.DB) contractUser.UserRepository {
			return repositoriesUser.NewUserRepository(dbadapter.NewGormAdapter(db))
		}),
		fx.Provide(func(db *gorm.DB) contractInventory.InventoryRepository {
			return repositoriesInventory.NewInventoryRepository(dbadapter.NewGormAdapter(db))
		}),

		fx.Provide(func(db *gorm.DB) contractProduct.ProductRepository {
			return repositoryProduct.NewProductRepository(dbadapter.NewGormAdapter(db))
		}),

		// Storage
		fx.Provide(getStorage),

		// Services
		fx.Provide(func(storage storage.Storage) qrcode.QRCodeService {
			return qrcode.NewQRCodeServiceWithStorage(storage)
		}),
		fx.Provide(func(storage storage.Storage) pallet.PalletExportService {
			return pallet.NewPalletExportService(storage)
		}),
		fx.Provide(func(repo contractPallet.PalletRepository, qrService qrcode.QRCodeService, palletRackRepo contractPalletRack.PalletRackRepository, storage storage.Storage, exportService pallet.PalletExportService) pallet.PalletService {
			return pallet.NewPalletService(repo, qrService, palletRackRepo, storage, exportService)
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

		fx.Provide(func(repo contractInventory.InventoryRepository) inventory.InventoryService {
			return inventory.NewInventoryService(repo)
		}),

		fx.Provide(func(repository contractProduct.ProductRepository) productService.ProductService {
			return productService.NewProductService(repository)
		}),

		// Controllers
		fx.Provide(palletController.NewPalletController),
		fx.Provide(palletizedProductController.NewPalletizedProductController),
		fx.Provide(palletRackController.NewPalletRackController),
		fx.Provide(palletRackController.NewAdminPalletRackController),
		fx.Provide(authController.NewAuthController),
		fx.Provide(userController.NewUserController),
		fx.Provide(inventoryController.NewInventoryController),
		fx.Provide(productController.NewProductController),

		// Middlewares
		fx.Provide(middlewares.NewAuthMiddleware),
		fx.Provide(middlewares.NewInventoryMiddleware),
	)
}

func getStorage() storage.Storage {
	// Decide provider via env var STORAGE_PROVIDER: "minio" or "local" (default)
	provider := os.Getenv("STORAGE_PROVIDER")
	if provider == "minio" {
		endpoint := os.Getenv("MINIO_ENDPOINT")
		accessKey := os.Getenv("MINIO_ACCESS_KEY")
		secretKey := os.Getenv("MINIO_SECRET_KEY")
		bucket := os.Getenv("MINIO_BUCKET")
		region := os.Getenv("MINIO_REGION")
		useSSL := os.Getenv("MINIO_USE_SSL") == "true"
		if endpoint == "" || accessKey == "" || secretKey == "" || bucket == "" {
			// fallback to local
			baseDir := os.Getenv("STORAGE_BASE_DIR")
			baseURL := os.Getenv("STORAGE_BASE_URL")
			if baseDir == "" {
				baseDir = "storage"
			}
			if baseURL == "" {
				baseURL = "http://localhost:3000/"
			}
			return storage.NewLocalStorage(baseDir, baseURL)
		}
		s3, err := storage.NewS3Storage(endpoint, accessKey, secretKey, bucket, region, useSSL)
		if err != nil {
			// fallback to local
			baseDir := os.Getenv("STORAGE_BASE_DIR")
			baseURL := os.Getenv("STORAGE_BASE_URL")
			if baseDir == "" {
				baseDir = "storage"
			}
			if baseURL == "" {
				baseURL = "http://localhost:3000/"
			}
			return storage.NewLocalStorage(baseDir, baseURL)
		}
		return s3
	}

	// default local storage
	baseDir := os.Getenv("STORAGE_BASE_DIR")
	baseURL := os.Getenv("STORAGE_BASE_URL")
	if baseDir == "" {
		baseDir = "storage"
	}
	if baseURL == "" {
		baseURL = "http://localhost:3000/"
	}
	return storage.NewLocalStorage(baseDir, baseURL)
}
