package main

import (
	"log"
	"net/http"
	"os"
	"time"

	routes "go_inventory/SupplyInventory/Application/Routes"
	db "go_inventory/SupplyInventory/Infrastructure/Db"

	"github.com/gin-contrib/cors"
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

	// Habilitar CORS para qualquer origem
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
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
			c.JSON(http.StatusNotFound, gin.H{"error": "rota não encontrada"})
			return
		}
		if len(path) >= 7 && path[:7] == "/qrcodes" {
			c.JSON(http.StatusNotFound, gin.H{"error": "arquivo não encontrado"})
			return
		}
		c.File("./frontend/build/index.html")
	})

	// Templates Go (se ainda precisar)
	router.LoadHTMLGlob("templates/**/*.tmpl")

	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum .env encontrado, usando variáveis do sistema")
	}

	// Conectar e migrar DB
	dbInstance := db.Connect()
	db.Migrate(dbInstance)

	// Registrar rotas API
	routes.RegisterRoutes(router, dbInstance)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	server := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		MaxHeaderBytes: 1 << 60, // 1MB, pode aumentar se precisar
	}

	server.ListenAndServe()
}
