package controllers

import (
	"net/http"

	services "go_inventory/SupplyInventory/Application/Services"
	domain "go_inventory/SupplyInventory/Domain"

	"github.com/gin-gonic/gin"
)

type PalletizedProductController struct {
	palletizedProductService services.PalletizedProductService
}

func NewPalletizedProductController(palletizedProductService services.PalletizedProductService) *PalletizedProductController {
	return &PalletizedProductController{palletizedProductService: palletizedProductService}
}

type PalletizedProductRequest struct {
	PalletizedProduct domain.PalletizedProductEntity `json:"PalletizedProduct" binding:"required"`
}

func (controller *PalletizedProductController) RegisterProductPallet(group *gin.RouterGroup) {
	group.PATCH("/:palletId", controller.addProductsToPallet)
}

func (controller *PalletizedProductController) addProductsToPallet(c *gin.Context) {
	var req PalletizedProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "O body deve conter um produto com o palete, ean e quantidade de produtos",
		})
		return
	}

	updatedPallet, err := controller.palletizedProductService.AddProductsToPallet(req.PalletizedProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPallet)
}
