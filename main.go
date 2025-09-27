package main

import (
	"log"
	"os"

	webRoutes "go_inventory/Front/Application/WebRoutes"
	routes "go_inventory/SupplyInventory/Application/Routes"
	db "go_inventory/SupplyInventory/Infrastructure/Db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go_inventory/docs"
)

// @title           Go Inventory API
// @version         1.0
// @description     It is system to manage pallets and related products.

// @contact.name   API Support
// @contact.email  joaoferrazp@gmail.com

// @host localhost:3000

func main() {
	router := gin.Default()

	// Pasta estática para QR codes
	router.Static("/qrcodes", "./storage/qrcodes")

	router.Static("/static", "./static")

	// Templates
	router.LoadHTMLGlob("templates/**/*.tmpl")

	// Carregar variáveis de ambiente

	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum .env encontrado, usando variáveis do sistema")
	}

	// Conectar e migrar DB
	dbInstance := db.Connect()
	db.Migrate(dbInstance)

	// Registrar rotas passando a instância do DB
	routes.RegisterRoutes(router, dbInstance)
	webRoutes.RegisterWebRoutes(router, dbInstance)

	// Pegar porta do .env
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":" + port)
}
