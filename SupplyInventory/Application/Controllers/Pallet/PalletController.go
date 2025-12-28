package controllers

import (
	"net/http"

	requestsHelper "go_inventory/Helpers/RequestsHelper"
	requests "go_inventory/SupplyInventory/Application/Requests"
	pallet "go_inventory/SupplyInventory/Application/Services/Pallet"

	"github.com/gin-gonic/gin"
)

type PalletController struct {
	service pallet.PalletService
}

func NewPalletController(service pallet.PalletService) *PalletController {
	return &PalletController{service: service}
}

func (controller *PalletController) Register(group *gin.RouterGroup) {
	group.GET("/", controller.ListPallets)
	group.GET("/:id", controller.FindPalletById)
	group.PATCH("/:id", controller.UpdatePallet)
	group.POST("/", controller.CreatePallet)
	group.DELETE("/:id", controller.DeletePalletById)
}

// @Summary List pallets
// @Tags Pallets
// @Accept json
// @Produce json
// @Success 200 {array} entities.PalletEntity
// @Failure 404 "Not Found"
// @Router /api/v1/pallets [get]
func (controller *PalletController) ListPallets(c *gin.Context) {
	pallets, appErr := controller.service.ListPallets()
	if appErr != nil {
		c.JSON(appErr.ErrorCode(), appErr.Error())
	}

	c.JSON(http.StatusOK, pallets)
}

// @Summary Get pallet by ID
// @Tags Pallets
// @Accept json
// @Produce json
// @Success 200 {object} entities.PalletEntity
// @Failure 404 "Not Found"
// @Param id path int true "ID do pallet"
// @Router /api/v1/pallets/{id} [get]
func (controller *PalletController) FindPalletById(c *gin.Context) {
	id, err := requestsHelper.GetIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
		return
	}

	pallet, appErr := controller.service.FindPalletById(id)
	if pallet == nil {
		c.JSON(appErr.ErrorCode(), gin.H{"message": appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, pallet)
}

// @Summary Delete Pallet
// @Tags Pallets
// @Accept json
// @Produce json
// @Param pallet body requests.PalletRequest true "Palletized Product"
// @Success 200 {object} entities.PalletEntity
// @Failure 422 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Param id path int true "ID do pallet"
// @Router /api/v1/pallets/{id} [delete]
func (controller *PalletController) DeletePalletById(c *gin.Context) {
	id, err := requestsHelper.GetIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "ID inválido"})
		return
	}

	result, appErr := controller.service.DeletePalletById(id)

	if appErr != nil {
		c.JSON(appErr.ErrorCode(), gin.H{"message": appErr.Error()})
		return
	}

	c.JSON(http.StatusNoContent, result)
}

// @Summary Create Pallet
// @Tags Pallets
// @Accept json
// @Produce json
// @Param pallet body requests.PalletRequest true "Palletized Product"
// @Success 200 {object} entities.PalletEntity
// @Failure 422 {object} map[string]string
// @Router /api/v1/pallets [post]
func (controller *PalletController) CreatePallet(c *gin.Context) {
	var req requests.PalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": requestsHelper.FormatValidationErrors(err)})
		return
	}

	newPallet, appErr := controller.service.CreatePallet(req.Name, req.PalletRackID)

	if appErr != nil {
		c.JSON(appErr.ErrorCode(), gin.H{"message": appErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, newPallet)
}

// @Summary Update Pallet
// @Tags Pallets
// @Accept json
// @Produce json
// @Param pallet body requests.PalletRequest true "Palletized Product"
// @Success 200 {object} entities.PalletEntity
// @Failure 422 {object} map[string]string
// @Param id path int true "ID do pallet"
// @Router /api/v1/pallets/{id} [patch]
func (controller *PalletController) UpdatePallet(c *gin.Context) {
	id, err := requestsHelper.GetIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "ID inválido"})
		return
	}

	var req requests.PalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	pallet, appErr := controller.service.UpdatePallet(id, req.Name, req.PalletRackID)
	if pallet == nil {
		c.JSON(appErr.ErrorCode(), gin.H{"message": appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, pallet)
}
