package main

import (
	"github.com/gin-gonic/gin"
	"go_inventory/SupplyInventory/Application/Routes"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, world! 12")
	})

	routes.RegisterRoutes(r)

	r.Run(":3000") // roda em http://localhost:3000
}
