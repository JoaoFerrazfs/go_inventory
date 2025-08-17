package controllers

import (
	"net/http"

	requestsHelper "go_inventory/Helpers/RequestsHelper"
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
		id, err := requestsHelper.GetIDParam(c, "id")
		if err != nil {
			c.JSON(400, gin.H{"message": "ID inválido"})
			return
		}

		position, err := services.FindPositionById(id)
		if position == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
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

		newPosition, err := services.CreatePosition(req.Position)
		if newPosition == nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, newPosition)
	})
}
