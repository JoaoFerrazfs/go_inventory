package main

import (
	routes "go_inventory/SupplyInventory/Application/Routes"
	db "go_inventory/SupplyInventory/Infrastructure/Db"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	db.Connect()
	db.Migrate()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, world! 12")
	})

	routes.RegisterRoutes(r)

	r.Run(":3000") // roda em http://localhost:3000
}
