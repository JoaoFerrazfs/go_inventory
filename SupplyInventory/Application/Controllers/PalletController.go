package controllers

import (
	"net/http"

	requestsHelper "go_inventory/Helpers/RequestsHelper"
	services "go_inventory/SupplyInventory/Application/Services"
	domain "go_inventory/SupplyInventory/Domain"

	"github.com/gin-gonic/gin"
)

type PalletRequest struct {
	Pallet domain.Pallet `json:"pallet" binding:"required"`
}

func Register(group *gin.RouterGroup) {
	group.GET("/", func(c *gin.Context) {
		palletPallets := services.ListPallets()

		c.JSON(http.StatusOK, palletPallets)
	})

	group.GET("/:id", func(c *gin.Context) {
		id, err := requestsHelper.GetIDParam(c, "id")
		if err != nil {
			c.JSON(400, gin.H{"message": "ID inválido"})
			return
		}

		pallet, err := services.FindPalletById(id)
		if pallet == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, pallet)
	})

	group.POST("/", func(c *gin.Context) {
		var req PalletRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		newPallet, err := services.CreatePallet(req.Pallet)
		if newPallet == nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, newPallet)
	})
}
