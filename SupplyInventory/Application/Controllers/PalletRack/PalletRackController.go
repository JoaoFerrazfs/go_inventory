package controllers

import (
	"net/http"

	development "go_inventory/Helpers/Development"
	requestsHelper "go_inventory/Helpers/RequestsHelper"
	requests "go_inventory/SupplyInventory/Application/Requests"
	services "go_inventory/SupplyInventory/Application/Services"

	"github.com/gin-gonic/gin"
)

type PalletRackController struct {
	palletRackService services.PalletRackService
}

func NewPalletRackController(palletRackService services.PalletRackService) *PalletRackController {
	return &PalletRackController{palletRackService: palletRackService}
}

func (controller *PalletRackController) RegisterPalletRack(group *gin.RouterGroup) {
	group.POST("/", controller.createPalletRack)
	group.GET("/", controller.listRacks)
	group.GET("/:id", controller.FindRackById)
	group.DELETE("/:id", controller.DeleteRack)
}

// @Summary Create Pallet Racks
// @Tags Pallet Racks
// @Accept json
// @Produce json
// @Param PalletRack body requests.PalletRackRequest true "Palletized Product"
// @Success 200 {object} entities.PalletRackEntity
// @Failure 422 {object} map[string]string
// @Router /api/v1/racks [post]
func (controller *PalletRackController) createPalletRack(c *gin.Context) {
	var req requests.PalletRackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		development.Dump(req)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": requestsHelper.FormatValidationErrors(err)})
		return
	}

	newPalletRack, err := controller.palletRackService.Create(req.Name, req.Location, req.TotalCapacity)

	if newPalletRack == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newPalletRack)
}

// @Summary List Pallet Racks
// @Tags Pallet Racks
// @Accept json
// @Produce json
// @Param PalletRack body requests.PalletRackRequest true "Palletized Product"
// @Success 200 {array} entities.PalletRackEntity
// @Failure 404 {object} map[string]string
// @Router /api/v1/racks [get]
func (controller *PalletRackController) listRacks(c *gin.Context) {
	racks, err := controller.palletRackService.ListRacks()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, racks)
}

// @Summary Get Rack
// @Tags Pallet Racks
// @Accept json
// @Produce json
// @Success 200 {object} entities.PalletRackEntity
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Param id path int true "ID do Rack"
// @Router /api/v1/racks/{id} [get]
func (pc *PalletRackController) FindRackById(c *gin.Context) {
	id, err := requestsHelper.GetIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "ID inválido"})
		return
	}

	pallet, appErr := pc.palletRackService.FindPalletById(id)
	if pallet == nil {
		c.JSON(appErr.ErrorCode(), gin.H{"message": appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, pallet)
}

// @Summary Delete Rack
// @Tags Pallet Racks
// @Accept json
// @Produce json
// @Success 204 "Delete successful"
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Param id path int true "ID do Rack"
// @Router /api/v1/racks/{id} [delete]
func (pc *PalletRackController) DeleteRack(c *gin.Context) {
	id, err := requestsHelper.GetIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "ID inválido"})
		return
	}

	_, appErr := pc.palletRackService.DeleteRack(id)
	if appErr != nil {
		c.JSON(appErr.ErrorCode(), gin.H{"message": appErr.Error()})
		return
	}

	c.Status(http.StatusOK)
}
