package main

import (
	"log"
	"os"

	routes "go_inventory/SupplyInventory/Application/Routes"
	db "go_inventory/SupplyInventory/Infrastructure/Db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	// Pasta estática para QR codes
	r.Static("/qrcodes", "./storage/qrcodes")

	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum .env encontrado, usando variáveis do sistema")
	}

	// Conectar e migrar DB
	dbInstance := db.Connect()
	db.Migrate(dbInstance)

	// Registrar rotas passando a instância do DB
	routes.RegisterRoutes(r, dbInstance)

	// Pegar porta do .env
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	r.Run(":" + port)
}
