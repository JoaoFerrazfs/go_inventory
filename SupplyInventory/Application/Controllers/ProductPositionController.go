package controllers

import (
	"log"
	"net/http"

	services "go_inventory/SupplyInventory/Application/Services"
	domain "go_inventory/SupplyInventory/Domain"

	"github.com/gin-gonic/gin"
)

type ProductsRequest struct {
	Product domain.PositionProduct `json:"product" binding:"required"`
}

func RegisterProductPosition(group *gin.RouterGroup) {
	group.PATCH("/:position/products", func(c *gin.Context) {
		var req ProductsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Default().Println("Error binding JSON:", req)
			c.JSON(400, gin.H{"message": "O body deve conter products como array de números"})
			return
		}

		updatedPosition := services.AddProductsToPositition(req.Product)

		c.JSON(http.StatusOK, updatedPosition)
	})
}
