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
	group.GET("/", pc.ListPallets)
	group.GET("/:id", pc.FindPalletById)
	group.POST("/", pc.CreatePallet)
	group.DELETE("/:id", pc.DeletePalletById)
}

// @Summary List pallets
// @Tags Pallets
// @Accept json
// @Produce json
// @Success 200 {array} domain.PalletEntity
// @Failure 404 "Not Found"
// @Router /api/v1/pallets [get]
func (pc *PalletController) ListPallets(c *gin.Context) {
	pallets, appErr := pc.service.ListPallets()
	if appErr != nil {
		c.JSON(appErr.ErrorCode(), appErr.Error())
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
func (pc *PalletController) FindPalletById(c *gin.Context) {
	id, err := requestsHelper.GetIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
		return
	}

	pallet, appErr := pc.service.FindPalletById(id)
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
// @Success 200 {object} domain.PalletEntity
// @Failure 422 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Param id path int true "ID do pallet"
// @Router /api/v1/pallets/{id} [delete]
func (pc *PalletController) DeletePalletById(c *gin.Context) {
	id, err := requestsHelper.GetIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "ID inválido"})
		return
	}

	result, appErr := pc.service.DeletePalletById(id)

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
// @Success 200 {object} domain.PalletEntity
// @Failure 422 {object} map[string]string
// @Router /api/v1/pallets [post]
func (pc *PalletController) CreatePallet(c *gin.Context) {
	var req requests.PalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	newPallet, appErr := pc.service.CreatePallet(req.Name, req.PalletRackID)

	if appErr != nil {
		c.JSON(appErr.ErrorCode(), gin.H{"message": appErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, newPallet)
}
