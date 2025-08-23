package controllers

import (
	"net/http"

	services "go_inventory/SupplyInventory/Application/Services"
	domain "go_inventory/SupplyInventory/Domain"

	"github.com/gin-gonic/gin"
)

type PalletRackRequest struct {
	PalletRack domain.PalletRackEntity `json:"palletRack" binding:"required"`
}

type PalletRackController struct {
	palletRackService services.PalletRackService
}

func NewPalletRackController(palletRackService services.PalletRackService) *PalletRackController {
	return &PalletRackController{palletRackService: palletRackService}
}

func (controller *PalletRackController) RegisterPalletRack(group *gin.RouterGroup) {
	group.POST("/", controller.createPalletRack)
}

func (controller *PalletRackController) createPalletRack(c *gin.Context) {
	var req PalletRackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	}

	newPalletRack, err := controller.palletRackService.Create(req.PalletRack)

	if newPalletRack == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newPalletRack)
}
