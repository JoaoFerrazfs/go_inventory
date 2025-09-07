package controllers

import (
	"net/http"

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
}

// @Summary Create Pallet Racks
// @Tags Pallet Racks
// @Accept json
// @Produce json
// @Param PalletRack body requests.PalletRackRequest true "Palletized Product"
// @Success 200 {object} domain.PalletRackEntity
// @Failure 422 {object} map[string]string
// @Router /api/v1/racks [post]
func (controller *PalletRackController) createPalletRack(c *gin.Context) {
	var req requests.PalletRackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	}

	newPalletRack, err := controller.palletRackService.Create(req.Name)

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
// @Success 200 {object} domain.PalletRackEntity
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
