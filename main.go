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

	r.Static("/qrcodes", "./storage/qrcodes")

	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum .env encontrado, usando variáveis do sistema")
	}

	db.Connect()
	db.Migrate()

	routes.RegisterRoutes(r)

	port := os.Getenv("PORT")
	r.Run(":" + port)
}
