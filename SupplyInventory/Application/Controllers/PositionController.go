package controllers

import (
	"net/http"
	"strconv"

	services "go_inventory/SupplyInventory/Application/Services"
	domain "go_inventory/SupplyInventory/Domain"

	"github.com/gin-gonic/gin"
)

type PositionRequest struct {
	Position domain.Position `json:"position" binding:"required"`
}

func Register(group *gin.RouterGroup) {
	group.GET("/", func(c *gin.Context) {
		positions := services.ListPositions()
		c.JSON(http.StatusOK, positions)
	})

	group.GET("/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		idUint64, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
			return
		}

		id := uint(idUint64)

		position := services.FindPositionById(id)
		if position == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Position not found"})
			return
		}
		c.JSON(http.StatusOK, position)
	})

	group.POST("/", func(c *gin.Context) {
		var req PositionRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		newPosition := services.CreatePosition(req.Position)
		if newPosition == nil {
			c.JSON(http.StatusBadGateway, gin.H{"message": "Position could not be created"})
			return
		}

		c.JSON(http.StatusCreated, newPosition)
	})
}
