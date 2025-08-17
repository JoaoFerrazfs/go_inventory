package controllers

import (
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
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "O body deve conter um produto com o palete, ean e quantidade de produtos "})
			return
		}

		updatedPosition := services.AddProductsToPositition(req.Product)

		c.JSON(http.StatusOK, updatedPosition)
	})
}
