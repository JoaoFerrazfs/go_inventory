package controllers

import (
	"net/http"

	services "go_inventory/SupplyInventory/Application/Services"
	domain "go_inventory/SupplyInventory/Domain"

	"github.com/gin-gonic/gin"
)

type PalletizedProductController struct {
	service services.PalletService
}

func NewPalletizedProductController(service services.PalletService) *PalletizedProductController {
	return &PalletizedProductController{service: service}
}

type PalletizedProductRequest struct {
	PalletizedProduct domain.PalletizedProductEntity `json:"PalletizedProduct" binding:"required"`
}

func (pc *PalletizedProductController) RegisterProductPallet(group *gin.RouterGroup) {
	group.PATCH("/:pallet", pc.addProductsToPallet)
}

func (pc *PalletizedProductController) addProductsToPallet(c *gin.Context) {
	var req PalletizedProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "O body deve conter um produto com o palete, ean e quantidade de produtos",
		})
		return
	}

	updatedPallet := pc.service.AddProductsToPallet(req.PalletizedProduct)
	if updatedPallet == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao adicionar produto ao palete"})
		return
	}

	c.JSON(http.StatusOK, updatedPallet)
}
