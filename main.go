package main

import (
	"github.com/gin-gonic/gin"
	"go_inventory/SupplyInventory/Application/Controllers"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, world!")
	})

	controllers.RegisterRoutes(r)

	r.Run(":3000") // roda em http://localhost:3000
}
