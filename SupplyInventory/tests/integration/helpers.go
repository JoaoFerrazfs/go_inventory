package integration

import (
	"context"
	container "go_inventory/Container"
	auth "go_inventory/SupplyInventory/Application/Controllers/Auth"
	inventory "go_inventory/SupplyInventory/Application/Controllers/Inventory"
	pallet "go_inventory/SupplyInventory/Application/Controllers/Pallet"
	palletRack "go_inventory/SupplyInventory/Application/Controllers/PalletRack"
	palletizedProduct "go_inventory/SupplyInventory/Application/Controllers/PalletizedProduct"
	user "go_inventory/SupplyInventory/Application/Controllers/User"
	middlewares "go_inventory/SupplyInventory/Application/Middlewares"
	inventoryService "go_inventory/SupplyInventory/Application/Services/Inventory"
	jwtService "go_inventory/SupplyInventory/Application/Services/Jwt"
	palletService "go_inventory/SupplyInventory/Application/Services/Pallet"
	palletRackService "go_inventory/SupplyInventory/Application/Services/PalletRack"
	palletizedProductService "go_inventory/SupplyInventory/Application/Services/PalletizedProduct"
	qrCodeService "go_inventory/SupplyInventory/Application/Services/QrCode"
	userService "go_inventory/SupplyInventory/Application/Services/User"
	entities "go_inventory/SupplyInventory/Domain/Entities" //nolint
	Inventory "go_inventory/SupplyInventory/Domain/contracts/repositories/Inventory"
	Pallet "go_inventory/SupplyInventory/Domain/contracts/repositories/Pallet"
	PalletRack "go_inventory/SupplyInventory/Domain/contracts/repositories/PalletRack"
	PalletizedProduct "go_inventory/SupplyInventory/Domain/contracts/repositories/PalletizedProduct"
	User "go_inventory/SupplyInventory/Domain/contracts/repositories/User"
	testutils "go_inventory/SupplyInventory/tests/testutils"

	// infra repos for building services bound to provided DB
	inventoryInfra "go_inventory/SupplyInventory/Infrastructure/repositories/Inventory"
	palletInfra "go_inventory/SupplyInventory/Infrastructure/repositories/Pallet"
	palletRackInfra "go_inventory/SupplyInventory/Infrastructure/repositories/PalletRack"
	palletizedProductInfra "go_inventory/SupplyInventory/Infrastructure/repositories/PalletizedProduct"
	userInfra "go_inventory/SupplyInventory/Infrastructure/repositories/User"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"go.uber.org/fx"
	"gorm.io/gorm"

	dbadapter "go_inventory/SupplyInventory/Infrastructure/repositories/db"
	mocks "go_inventory/SupplyInventory/tests/mocks"
)

// TestDependencies holds all dependencies for integration tests
type TestDependencies struct {
	UserRepo                 User.UserRepository
	PalletRepo               Pallet.PalletRepository
	PalletizedProductRepo    PalletizedProduct.PalletizedProductRepository
	PalletRackRepo           PalletRack.PalletRackRepository
	InventoryRepo            Inventory.InventoryRepository
	UserService              userService.UserService
	PalletService            palletService.PalletService
	PalletizedProductService palletizedProductService.PalletizedProductService
	PalletRackService        palletRackService.PalletRackService
	InventoryService         inventoryService.InventoryService
	QrCodeService            qrCodeService.QRCodeService
	JwtService               jwtService.JWTService
}

// IntegrationTestHelper provides utilities for integration testing
// IMPORTANT: When adding new controllers/routes, create corresponding SetupRouterFor{Controller} method
type IntegrationTestHelper struct {
	DB *gorm.DB
	TestDependencies
	App *fx.App
}

func NewIntegrationTestHelper() *IntegrationTestHelper {
	gin.SetMode(gin.TestMode)
	db := testutils.SetupTestDB()

	var deps TestDependencies
	app := fx.New(
		container.BuildOptions(db),
		fx.Invoke(func(
			userRepo User.UserRepository,
			palletRepo Pallet.PalletRepository,
			palletizedProductRepo PalletizedProduct.PalletizedProductRepository,
			palletRackRepo PalletRack.PalletRackRepository,
			inventoryRepo Inventory.InventoryRepository,
			userSrv userService.UserService,
			palletSrv palletService.PalletService,
			palletizedProductSrv palletizedProductService.PalletizedProductService,
			palletRackSrv palletRackService.PalletRackService,
			inventorySrv inventoryService.InventoryService,
			jwtSrv jwtService.JWTService,
			qrSrv qrCodeService.QRCodeService,
		) {
			deps = TestDependencies{
				UserRepo:                 userRepo,
				PalletRepo:               palletRepo,
				PalletizedProductRepo:    palletizedProductRepo,
				PalletRackRepo:           palletRackRepo,
				InventoryRepo:            inventoryRepo,
				UserService:              userSrv,
				PalletService:            palletSrv,
				PalletizedProductService: palletizedProductSrv,
				PalletRackService:        palletRackSrv,
				InventoryService:         inventorySrv,
				JwtService:               jwtSrv,
				QrCodeService:            qrSrv,
			}
		}),
	)
	// Start the fx app so constructors run and Populate is executed
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}

	return &IntegrationTestHelper{
		DB:               db,
		TestDependencies: deps,
		App:              app,
	}
}

// Stop stops the underlying fx app to release resources used in tests
func (h *IntegrationTestHelper) Stop() error {
	if h.App == nil {
		return nil
	}
	return h.App.Stop(context.Background())
}

// Router Setup Methods
// Each controller must have a corresponding SetupRouterFor{Controller} method
// to ensure integration tests use the same dependency injection as production
func (h *IntegrationTestHelper) SetupRouterForAuth(db *gorm.DB) *gin.Engine {
	// Build user service bound to provided DB/tx so it sees uncommitted data in the same transaction
	userRepo := userInfra.NewUserRepository(dbadapter.NewGormAdapter(db))
	userSrv := userService.NewUserService(userRepo)
	controller := auth.NewAuthController(h.JwtService, userSrv)
	r := gin.Default()
	api := r.Group("/api/v1/auth")
	controller.RegisterLogin(api)
	return r
}

func (h *IntegrationTestHelper) SetupRouterForInventory(db *gorm.DB) *gin.Engine {
	inventoryRepo := inventoryInfra.NewInventoryRepository(dbadapter.NewGormAdapter(db))
	inventorySrv := inventoryService.NewInventoryService(inventoryRepo)
	userRepo := userInfra.NewUserRepository(dbadapter.NewGormAdapter(db))
	userSrv := userService.NewUserService(userRepo)

	controller := inventory.NewInventoryController(inventorySrv, userSrv)

	r := gin.Default()
	// Middleware to simulate user authentication for testing
	r.Use(func(c *gin.Context) {
		if userIDStr := c.GetHeader("X-Test-User-ID"); userIDStr != "" {
			if userID, err := strconv.Atoi(userIDStr); err == nil {
				c.Set("userID", uint(userID))
			}
		}
		c.Next()
	})
	api := r.Group("/api/v1/inventories")
	controller.Register(api)
	return r
}

func (h *IntegrationTestHelper) SetupRouterForPallet(db *gorm.DB) *gin.Engine {
	palletRepo := palletInfra.NewPalletRepository(dbadapter.NewGormAdapter(db))
	palletRackRepo := palletRackInfra.NewPalletRackRepository(dbadapter.NewGormAdapter(db))
	qrSrv := h.QrCodeService
	// For tests, use mock storage to avoid creating real files
	mockStorage := &mocks.MockStorage{}
	mockStorage.On("Upload", mock.Anything, mock.Anything).Return("http://localhost:3000/reports/Pallets/test.csv", nil)
	mockStorage.On("GetURL", mock.Anything).Return("http://localhost:3000/reports/Pallets/test.csv", nil)
	exportSrv := palletService.NewPalletExportService(mockStorage)
	palletSrv := palletService.NewPalletService(palletRepo, qrSrv, palletRackRepo, mockStorage, exportSrv)
	controller := pallet.NewPalletController(palletSrv)
	r := gin.Default()
	api := r.Group("/api/v1/pallets")
	// inventory middleware
	invRepo := inventoryInfra.NewInventoryRepository(dbadapter.NewGormAdapter(db))
	invMiddleware := middlewares.NewInventoryMiddleware(invRepo)
	api.Use(invMiddleware.Handler())
	controller.Register(api)
	return r
}

func (h *IntegrationTestHelper) SetupRouterForPalletizedProduct(db *gorm.DB) *gin.Engine {
	palletRepo := palletInfra.NewPalletRepository(dbadapter.NewGormAdapter(db))
	palletizedProductRepo := palletizedProductInfra.NewPalletizedProductRepository(dbadapter.NewGormAdapter(db), palletRepo)
	palletizedProductSrv := palletizedProductService.NewPalletizedProductService(palletRepo, palletizedProductRepo)
	controller := palletizedProduct.NewPalletizedProductController(palletizedProductSrv)
	r := gin.Default()
	api := r.Group("/api/v1/pallet/products")
	invRepo := inventoryInfra.NewInventoryRepository(dbadapter.NewGormAdapter(db))
	invMiddleware := middlewares.NewInventoryMiddleware(invRepo)
	api.Use(invMiddleware.Handler())
	controller.RegisterProductPallet(api)
	return r
}

func (h *IntegrationTestHelper) SetupRouterForAdminPalletRack(db *gorm.DB) *gin.Engine {
	palletRackRepo := palletRackInfra.NewPalletRackRepository(dbadapter.NewGormAdapter(db))
	palletRackSrv := palletRackService.NewPalletRackService(palletRackRepo)
	controller := palletRack.NewAdminPalletRackController(palletRackSrv)
	r := gin.Default()
	// Add auth middleware simulation
	r.Use(func(c *gin.Context) {
		// Simulate authenticated user
		c.Set("userID", uint(1))
		c.Next()
	})
	api := r.Group("/api/v1/admin/racks")
	controller.RegisterAdminPalletRack(api)
	return r
}

func (h *IntegrationTestHelper) SetupRouterForPalletRack(db *gorm.DB) *gin.Engine {
	palletRackRepo := palletRackInfra.NewPalletRackRepository(dbadapter.NewGormAdapter(db))
	palletRackSrv := palletRackService.NewPalletRackService(palletRackRepo)
	controller := palletRack.NewPalletRackController(palletRackSrv)
	r := gin.Default()
	api := r.Group("/api/v1/racks")
	// inventory middleware
	invRepo := inventoryInfra.NewInventoryRepository(dbadapter.NewGormAdapter(db))
	invMiddleware := middlewares.NewInventoryMiddleware(invRepo)
	api.Use(invMiddleware.Handler())
	controller.RegisterPalletRack(api)
	return r
}

func (h *IntegrationTestHelper) SetupRouterForUser(db *gorm.DB) *gin.Engine {
	userRepo := userInfra.NewUserRepository(dbadapter.NewGormAdapter(db))
	userSrv := userService.NewUserService(userRepo)
	controller := user.NewUserController(userSrv)
	r := gin.Default()
	api := r.Group("/api/v1/users")
	controller.RegisterUserRoutes(api)
	return r
}

// SetupTestRouter creates a complete router with all routes configured for testing
// This centralizes route setup and avoids duplication with main.go
func SetupTestRouter(db *gorm.DB) *gin.Engine {
	var deps TestDependencies
	app := fx.New(
		container.BuildOptions(db),
		fx.Populate(&deps),
	)
	defer app.Stop(context.Background())

	// Initialize controllers
	authCtrl := auth.NewAuthController(deps.JwtService, deps.UserService)
	inventoryCtrl := inventory.NewInventoryController(deps.InventoryService, deps.UserService)
	palletCtrl := pallet.NewPalletController(deps.PalletService)
	palletizedProductCtrl := palletizedProduct.NewPalletizedProductController(deps.PalletizedProductService)
	palletRackCtrl := palletRack.NewPalletRackController(deps.PalletRackService)
	userCtrl := user.NewUserController(deps.UserService)

	// Setup router (similar to main.go setupRouter)
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Register routes (similar to main.go registerRoutes)
	apiV1 := router.Group("/api/v1")

	inventoriesGroup := apiV1.Group("/inventories")
	inventoryCtrl.Register(inventoriesGroup)

	palletsGroup := apiV1.Group("/pallets")
	palletCtrl.Register(palletsGroup)

	palletProductsGroup := apiV1.Group("/pallet/products")
	palletizedProductCtrl.RegisterProductPallet(palletProductsGroup)

	racksGroup := apiV1.Group("/racks")
	palletRackCtrl.RegisterPalletRack(racksGroup)

	authCtrl.RegisterLogin(apiV1.Group("/auth"))
	userCtrl.RegisterUserRoutes(apiV1.Group("/users"))

	return router
}

// Fixtures
func (h *IntegrationTestHelper) TruncateTables(db *gorm.DB) {
	testutils.TruncateTables(db)
}

func (h *IntegrationTestHelper) CreateTestUser(db *gorm.DB, name, email, password string) *entities.UserEntity {
	return testutils.CreateTestUser(db, name, email, password)
}

func (h *IntegrationTestHelper) CreateTestInventory(db *gorm.DB) *entities.InventoryEntity {
	return testutils.CreateTestInventory(db)
}

func (h *IntegrationTestHelper) CreateTestPallet(db *gorm.DB, name string, palletRackID uint) *entities.PalletEntity {
	return testutils.CreateTestPallet(db, name, palletRackID)
}

func (h *IntegrationTestHelper) CreateTestPalletRack(db *gorm.DB, name, location string, totalCapacity int) *entities.PalletRackEntity {
	return testutils.CreateTestPalletRack(db, name, location, totalCapacity)
}

func (h *IntegrationTestHelper) CreateTestProduct(db *gorm.DB, ean, name string) *entities.ProductEntity {
	return testutils.CreateTestProduct(db, ean, name)
}
