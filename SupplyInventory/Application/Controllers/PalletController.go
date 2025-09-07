package controllers

import (
	"net/http"

	requestsHelper "go_inventory/Helpers/RequestsHelper"
	requests "go_inventory/SupplyInventory/Application/Requests"
	services "go_inventory/SupplyInventory/Application/Services"

	"github.com/gin-gonic/gin"
)

type PalletController struct {
	service services.PalletService
}

func NewPalletController(service services.PalletService) *PalletController {
	return &PalletController{service: service}
}

func (pc *PalletController) Register(group *gin.RouterGroup) {
	group.GET("/", pc.listPallets)
	group.GET("/:id", pc.findPalletById)
	group.POST("/", pc.createPallet)
}

// @Summary List pallets
// @Tags Pallets
// @Accept json
// @Produce json
// @Success 200 {array} domain.PalletEntity
// @Failure 404 "Not Found"
// @Router /api/v1/pallets [get]
func (pc *PalletController) listPallets(c *gin.Context) {
	pallets, err := pc.service.ListPallets()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
	}

	c.JSON(http.StatusOK, pallets)
}

// @Summary Get pallet by ID
// @Tags Pallets
// @Accept json
// @Produce json
// @Success 200 {object} domain.PalletEntity
// @Failure 404 "Not Found"
// @Param id path int true "ID do pallet"
// @Router /api/v1/pallets/{id} [get]
func (pc *PalletController) findPalletById(c *gin.Context) {
	id, err := requestsHelper.GetIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
		return
	}

	pallet, err := pc.service.FindPalletById(id)
	if pallet == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pallet)
}

// @Summary Create Pallet
// @Tags Pallets
// @Accept json
// @Produce json
// @Param pallet body requests.PalletRequest true "Palletized Product"
// @Success 200 {object} domain.PalletEntity
// @Failure 422 {object} map[string]string
// @Router /api/v1/pallets [post]
func (pc *PalletController) createPallet(c *gin.Context) {
	var req requests.PalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	newPallet, err := pc.service.CreatePallet(req.Name, req.PalletRackID)

	if newPallet == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newPallet)
}
