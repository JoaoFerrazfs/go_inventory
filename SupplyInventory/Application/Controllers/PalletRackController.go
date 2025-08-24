package controllers

import (
	"net/http"

	services "go_inventory/SupplyInventory/Application/Services"

	"github.com/gin-gonic/gin"
)

type PalletRackRequest struct {
	PalletRack struct {
		Name string `json:"name"`
	} `json:"palletRack"`
}

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
// @Param PalletRack body PalletRackRequest true "Palletized Product"
// @Success 200 {object} domain.PalletRackEntity
// @Failure 422 {object} map[string]string
// @Router /api/v1/racks [post]
func (controller *PalletRackController) createPalletRack(c *gin.Context) {
	var req PalletRackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	}

	newPalletRack, err := controller.palletRackService.Create(req.PalletRack.Name)

	if newPalletRack == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newPalletRack)
}

func (controller *PalletRackController) listRacks(c *gin.Context) {
	racks, err := controller.palletRackService.ListRacks()

	if racks == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, racks)
}
