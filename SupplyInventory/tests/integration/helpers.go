package integration

import (
	"context"
	container "go_inventory/Container"
	auth "go_inventory/SupplyInventory/Application/Controllers/Auth"
	pallet "go_inventory/SupplyInventory/Application/Controllers/Pallet"
	palletRack "go_inventory/SupplyInventory/Application/Controllers/PalletRack"
	palletizedProduct "go_inventory/SupplyInventory/Application/Controllers/PalletizedProduct"
	user "go_inventory/SupplyInventory/Application/Controllers/User"
	jwtService "go_inventory/SupplyInventory/Application/Services/Jwt"
	palletService "go_inventory/SupplyInventory/Application/Services/Pallet"
	palletRackService "go_inventory/SupplyInventory/Application/Services/PalletRack"
	palletizedProductService "go_inventory/SupplyInventory/Application/Services/PalletizedProduct"
	qrCodeService "go_inventory/SupplyInventory/Application/Services/QrCode"
	userService "go_inventory/SupplyInventory/Application/Services/User"
	entities "go_inventory/SupplyInventory/Domain/Entities" //nolint
	Pallet "go_inventory/SupplyInventory/Domain/contracts/repositories/Pallet"
	PalletRack "go_inventory/SupplyInventory/Domain/contracts/repositories/PalletRack"
	PalletizedProduct "go_inventory/SupplyInventory/Domain/contracts/repositories/PalletizedProduct"
	User "go_inventory/SupplyInventory/Domain/contracts/repositories/User"
	testutils "go_inventory/SupplyInventory/tests/testutils"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// TestDependencies holds all dependencies for integration tests
type TestDependencies struct {
	fx.Out

	UserRepo                 User.UserRepository
	PalletRepo               Pallet.PalletRepository
	PalletizedProductRepo    PalletizedProduct.PalletizedProductRepository
	PalletRackRepo           PalletRack.PalletRackRepository
	UserService              userService.UserService
	PalletService            palletService.PalletService
	PalletizedProductService palletizedProductService.PalletizedProductService
	PalletRackService        palletRackService.PalletRackService
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
		fx.Populate(&deps),
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
	controller := auth.NewAuthController(h.JwtService, h.UserService)
	r := gin.Default()
	api := r.Group("/api/v1/auth")
	controller.RegisterLogin(api)
	return r
}

func (h *IntegrationTestHelper) SetupRouterForPallet(db *gorm.DB) *gin.Engine {
	controller := pallet.NewPalletController(h.PalletService)
	r := gin.Default()
	api := r.Group("/api/v1/pallets")
	controller.Register(api)
	return r
}

func (h *IntegrationTestHelper) SetupRouterForPalletizedProduct(db *gorm.DB) *gin.Engine {
	controller := palletizedProduct.NewPalletizedProductController(h.PalletizedProductService)
	r := gin.Default()
	api := r.Group("/api/v1/palletized-products")
	controller.RegisterProductPallet(api)
	return r
}

func (h *IntegrationTestHelper) SetupRouterForPalletRack(db *gorm.DB) *gin.Engine {
	controller := palletRack.NewPalletRackController(h.PalletRackService)
	r := gin.Default()
	api := r.Group("/api/v1/pallet-racks")
	controller.RegisterPalletRack(api)
	return r
}

func (h *IntegrationTestHelper) SetupRouterForUser(db *gorm.DB) *gin.Engine {
	controller := user.NewUserController(h.UserService)
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

	palletsGroup := apiV1.Group("/pallets")
	palletCtrl.Register(palletsGroup)

	palletProductsGroup := apiV1.Group("/palletized-products")
	palletizedProductCtrl.RegisterProductPallet(palletProductsGroup)

	racksGroup := apiV1.Group("/pallet-racks")
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

func (h *IntegrationTestHelper) CreateTestPallet(db *gorm.DB, name string, palletRackID uint) *entities.PalletEntity {
	return testutils.CreateTestPallet(db, name, palletRackID)
}

func (h *IntegrationTestHelper) CreateTestPalletRack(db *gorm.DB, name, location string, totalCapacity int) *entities.PalletRackEntity {
	return testutils.CreateTestPalletRack(db, name, location, totalCapacity)
}
