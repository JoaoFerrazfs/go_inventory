package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	container "go_inventory/Container"
	db "go_inventory/SupplyInventory/Infrastructure/Db"

	authControllerPkg "go_inventory/SupplyInventory/Application/Controllers/Auth"
	inventoryControllerPkg "go_inventory/SupplyInventory/Application/Controllers/Inventory"
	palletControllerPkg "go_inventory/SupplyInventory/Application/Controllers/Pallet"
	palletRackControllerPkg "go_inventory/SupplyInventory/Application/Controllers/PalletRack"
	palletizedProductControllerPkg "go_inventory/SupplyInventory/Application/Controllers/PalletizedProduct"
	productControllerPkg "go_inventory/SupplyInventory/Application/Controllers/Product"
	userControllerPkg "go_inventory/SupplyInventory/Application/Controllers/User"
	middlewares "go_inventory/SupplyInventory/Application/Middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"gorm.io/gorm"

	_ "go_inventory/docs"
)

// @title           Go Inventory API
// @version         1.0
// @description     It is system to manage pallets and related products.

// @contact.name   API Support
// @contact.email  joaoferrazp@gmail.com

// @host localhost:8000

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Digite: "Bearer {token}" para se autenticar.

func main() {
	// Carregar variáveis de ambiente
	envFile := os.Getenv("ENV_FILE")
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			log.Printf("Não foi possível carregar %s: %v", envFile, err)
		} else {
			log.Printf("Arquivo de env carregado: %s", envFile)
		}
	} else {
		if err := godotenv.Load(); err != nil {
			log.Println("Nenhum .env encontrado, usando variáveis do sistema")
		}
	}

	// Conectar DB
	dbInstance := db.Connect()

	// Criar router
	router := setupRouter()

	app := fx.New(
		container.BuildOptions(dbInstance),
		fx.Provide(func() *gin.Engine { return router }),
		fx.Invoke(migrateDB),
		fx.Invoke(registerRoutes),
		fx.Invoke(startServer),
	)

	app.Run()
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	// Disable automatic redirect for trailing slash mismatches to avoid 307
	// redirects on OPTIONS preflight requests which can break CORS flows.
	router.RedirectTrailingSlash = true

	// Enable CORS for any origin.
	// When AllowCredentials is true, using "*" in AllowOrigins is not allowed,
	// so use AllowOriginFunc that returns true for any origin.
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "x-inventory-id"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Pasta estática para QR codes
	router.Static("/qrcodes", "./storage/qrcodes")

	// Servir arquivos estáticos do React (CSS, JS, imagens)
	router.Static("/static", "./frontend/build/static")

	// Redirecionar qualquer rota que não seja API ou QR code para index.html
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		if len(path) >= 4 && path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "rota não encontrada", "path": path})
			return
		}
		if len(path) >= 7 && path[:7] == "/qrcodes" {
			c.JSON(http.StatusNotFound, gin.H{"error": "arquivo não encontrado"})
			return
		}
		c.File("./frontend/build/index.html")
	})

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func migrateDB(dbInstance *gorm.DB) {
	log.Println("Iniciando migration do banco de dados...")
	if err := db.Migrate(dbInstance); err != nil {
		log.Fatalf("Erro ao rodar migration: %v", err)
	}
	log.Println("Migration concluída com sucesso.")
}

func registerRoutes(
	router *gin.Engine,
	palletController *palletControllerPkg.PalletController,
	palletizedProductController *palletizedProductControllerPkg.PalletizedProductController,
	palletRackController *palletRackControllerPkg.PalletRackController,
	adminPalletRackController *palletRackControllerPkg.AdminPalletRackController,
	productController *productControllerPkg.ProductController,
	authController *authControllerPkg.AuthController,
	userController *userControllerPkg.UserController,
	authMiddleware *middlewares.AuthMiddleware,
	rbacMiddleware *middlewares.RBACMiddleware,
	inventoryController *inventoryControllerPkg.InventoryController,
	inventoryMiddleware *middlewares.InventoryMiddleware,
) {
	apiV1 := router.Group("/api/v1")

	palletsGroup := apiV1.Group("/pallets")
	palletsGroup.Use(authMiddleware.Handler())
	palletsGroup.Use(inventoryMiddleware.Handler())
	palletsGroup.Use(rbacMiddleware.RequireAny())
	palletController.Register(palletsGroup)

	palletProductsGroup := apiV1.Group("/pallet/products")
	palletProductsGroup.Use(authMiddleware.Handler())
	palletProductsGroup.Use(inventoryMiddleware.Handler())
	palletProductsGroup.Use(rbacMiddleware.RequireAny())
	palletizedProductController.RegisterProductPallet(palletProductsGroup)

	racksGroup := apiV1.Group("/racks")
	racksGroup.Use(authMiddleware.Handler())
	racksGroup.Use(inventoryMiddleware.Handler())
	racksGroup.Use(rbacMiddleware.RequireAny())
	palletRackController.RegisterPalletRack(racksGroup)

	adminRacksGroup := apiV1.Group("/admin/racks")
	adminRacksGroup.Use(authMiddleware.Handler())
	adminRacksGroup.Use(rbacMiddleware.RequireAdmin())
	adminPalletRackController.RegisterAdminPalletRack(adminRacksGroup)

	authController.RegisterLogin(apiV1.Group("/auth"))
	userController.RegisterUserRoutes(apiV1.Group("/users"))

	inventoryGroup := apiV1.Group("/inventories")
	inventoryGroup.Use(authMiddleware.Handler())
	inventoryGroup.Use(rbacMiddleware.RequireAny())
	inventoryController.Register(inventoryGroup)

	productGroup := apiV1.Group("/products")
	productGroup.Use(authMiddleware.Handler())
	productGroup.Use(rbacMiddleware.RequireAny())
	productController.Register(productGroup)
}

func startServer(lc fx.Lifecycle, router *gin.Engine) {
	port := os.Getenv("PORT")
	server := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		MaxHeaderBytes: 1 << 60,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting server...")
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping server...")
			return server.Shutdown(ctx)
		},
	})
}
