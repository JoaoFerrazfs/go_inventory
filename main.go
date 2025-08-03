package main

import (
	"github.com/gin-gonic/gin"
	"go_inventory/SupplyInventory/Application/Routes"
	"go_inventory/SupplyInventory/Infrastructure/Db"
)

func main() {
	r := gin.Default()
	db.Connect()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, world! 12")
	})

	routes.RegisterRoutes(r)

	r.Run(":3000") // roda em http://localhost:3000
}
