package controllers

import (
	"net/http"

	services "go_inventory/SupplyInventory/Application/Services"
	domain "go_inventory/SupplyInventory/Domain"

	"github.com/gin-gonic/gin"
)

type PalletizedProductRequest struct {
	PalletizedProduct domain.PalletizedProductEntity `json:"PalletizedProduct" binding:"required"`
}

func RegisterProductPallet(group *gin.RouterGroup) {
	group.PATCH("/:pallet/products", func(c *gin.Context) {
		var req PalletizedProductRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "O body deve conter um produto com o palete, ean e quantidade de produtos "})
			return
		}

		updatedPallet := services.AddProductsToPallet(req.PalletizedProduct)

		c.JSON(http.StatusOK, updatedPallet)
	})
}
