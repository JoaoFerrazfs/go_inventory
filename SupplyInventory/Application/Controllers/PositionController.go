package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"go_inventory/SupplyInventory/Application/Services"
	"go_inventory/SupplyInventory/Domain"
)

func RegisterRoutes(router *gin.Engine) {
	group := router.Group("/positions")

	group.GET("/", func(c *gin.Context) {
		positions := services.ListPositions()
		c.JSON(http.StatusOK, positions)
	})
	
	group.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		position := services.FindPositionById(id)
		if position == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Position not found"})
			return;
		}
		c.JSON(http.StatusOK, position)
	})

	group.POST("/", func(c *gin.Context) {
		var position domain.Position

		if err := c.BindJSON(&position); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		services.CreatePosition(position)
		c.JSON(http.StatusCreated, position)
	})
}