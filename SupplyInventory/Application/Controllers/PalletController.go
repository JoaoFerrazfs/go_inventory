package controllers

import (
	"net/http"

	requestsHelper "go_inventory/Helpers/RequestsHelper"
	services "go_inventory/SupplyInventory/Application/Services"
	domain "go_inventory/SupplyInventory/Domain"

	"github.com/gin-gonic/gin"
)

// Struct do controller
type PalletController struct {
	service services.PalletService
}

// Construtor
func NewPalletController(service services.PalletService) *PalletController {
	return &PalletController{service: service}
}

type PalletRequest struct {
	Pallet domain.Pallet `json:"pallet" binding:"required"`
}

// Método para registrar rotas
func (pc *PalletController) Register(group *gin.RouterGroup) {
	group.GET("/", pc.listPallets)
	group.GET("/:id", pc.findPalletById)
	group.POST("/", pc.createPallet)
}

// Handlers internos

func (pc *PalletController) listPallets(c *gin.Context) {
	pallets := pc.service.ListPallets()
	c.JSON(http.StatusOK, pallets)
}

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

func (pc *PalletController) createPallet(c *gin.Context) {
	var req PalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	newPallet, err := pc.service.CreatePallet(req.Pallet)
	if newPallet == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newPallet)
}
