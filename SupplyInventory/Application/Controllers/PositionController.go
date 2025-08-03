package controllers

import (
	"net/http"

	services "go_inventory/SupplyInventory/Application/Services"
	domain "go_inventory/SupplyInventory/Domain"

	"github.com/gin-gonic/gin"
)

func Register(group *gin.RouterGroup) {
	group.GET("/", func(c *gin.Context) {
		positions := services.ListPositions()
		c.JSON(http.StatusOK, positions)
	})

	group.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		position := services.FindPositionById(id)
		if position == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Position not found"})
			return
		}
		c.JSON(http.StatusOK, position)
	})

	group.POST("/", func(c *gin.Context) {
		var position domain.Position

		if err := c.ShouldBindJSON(&position); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		newPosition := services.CreatePosition(position)
		if newPosition == nil {
			c.JSON(http.StatusBadGateway, gin.H{"message": "Position could not be created"})
			return
		}

		c.JSON(http.StatusCreated, newPosition)
	})
}
