package controllers

import (
	"net/http"

	requestsHelper "go_inventory/Helpers/RequestsHelper"
	requests "go_inventory/SupplyInventory/Application/Requests"
	services "go_inventory/SupplyInventory/Application/Services"

	"github.com/gin-gonic/gin"
)

type PalletizedProductController struct {
	palletizedProductService services.PalletizedProductService
}

func NewPalletizedProductController(palletizedProductService services.PalletizedProductService) *PalletizedProductController {
	return &PalletizedProductController{palletizedProductService: palletizedProductService}
}

func (controller *PalletizedProductController) RegisterProductPallet(group *gin.RouterGroup) {
	group.PATCH("/:palletId", controller.addProductsToPallet)
	group.DELETE("/:palletId/:productsEan", controller.deleteProductsFromPallet)
}

// @Summary Add a product to a pallet
// @Tags Palletized products
// @Accept json
// @Produce json
// @Success 200 {object} domain.PalletizedProductEntity
// @Failure 422 {object} map[string]string
// @Param palletId path int true "ID do pallet"
// @Param pallet body requests.PalletizedProductRequest true "Palletized Product"
// @Router /api/v1/pallet/products/{palletId} [patch]
func (controller *PalletizedProductController) addProductsToPallet(c *gin.Context) {
	var req requests.PalletizedProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "O body deve conter um produto com o palete, ean e quantidade de produtos",
		})
		return
	}

	palletId, err := requestsHelper.GetIDParam(c, "palletId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "PalletId inválido"})
		return
	}

	updatedPallet, err := controller.palletizedProductService.AddProductsToPallet(palletId, req.EAN, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPallet)
}

// @Summary Delete product from pallet
// @Tags Palletized products
// @Accept json
// @Produce json
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Param palletId path int true "ID do pallet"
// @Param productsEan path int true "EAN do produto"
// @Router /api/v1/pallet/products/{palletId}/{productsEan} [Delete]
func (controller *PalletizedProductController) deleteProductsFromPallet(c *gin.Context) {
	palletId, err := requestsHelper.GetIDParam(c, "palletId")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Invalid palletId"})
		return
	}

	productsEan, err := requestsHelper.GetParamAsInt(c, "productsEan")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Invalid productsEan"})
		return
	}

	_, err = controller.palletizedProductService.DeleteProductsFromPallet(palletId, productsEan)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
