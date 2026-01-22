package controllers

import (
	"net/http"

	requestsHelper "go_inventory/Helpers/RequestsHelper"
	middlewares "go_inventory/SupplyInventory/Application/Middlewares"
	palletizedProductRequests "go_inventory/SupplyInventory/Application/Requests/PalletizedProduct"
	palletizedproduct "go_inventory/SupplyInventory/Application/Services/PalletizedProduct"

	"github.com/gin-gonic/gin"
)

type PalletizedProductController struct {
	palletizedProductService palletizedproduct.PalletizedProductService
}

func NewPalletizedProductController(palletizedProductService palletizedproduct.PalletizedProductService) *PalletizedProductController {
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
// @Security BearerAuth
// @Success 200 {object} entities.PalletizedProductEntity
// @Failure 422 {object} map[string]string
// @Param palletId path int true "ID do pallet"
// @Param pallet body palletizedproduct.PalletizedProductRequest true "Palletized Product"
// @Router /api/v1/pallet/products/{palletId} [patch]
func (controller *PalletizedProductController) addProductsToPallet(c *gin.Context) {
	var req palletizedProductRequests.PalletizedProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": requestsHelper.FormatValidationErrors(err)})
		return
	}

	palletId, err := requestsHelper.GetIDParam(c, "palletId")
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	inventoryID := middlewares.GetInventoryID(c)
	updatedPallet, appErr := controller.palletizedProductService.AddProductsToPallet(palletId, req.EAN, req.Quantity, inventoryID)

	if appErr != nil {
		c.JSON(appErr.ErrorCode(), gin.H{"error": appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPallet)
}

// @Summary Delete product from pallet
// @Tags Palletized products
// @Accept json
// @Produce json
// @Security BearerAuth
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

	_, appErr := controller.palletizedProductService.DeleteProductsFromPallet(palletId, productsEan)
	if appErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": appErr.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
