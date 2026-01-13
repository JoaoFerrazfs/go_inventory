package controllers

import (
	"net/http"
	"strconv"

	palletrack "go_inventory/SupplyInventory/Application/Services/PalletRack"

	"github.com/gin-gonic/gin"
)

type AdminPalletRackController struct {
	palletRackService palletrack.PalletRackService
}

func NewAdminPalletRackController(palletRackService palletrack.PalletRackService) *AdminPalletRackController {
	return &AdminPalletRackController{palletRackService: palletRackService}
}

func (controller *AdminPalletRackController) RegisterAdminPalletRack(group *gin.RouterGroup) {
	group.GET("", controller.ListRacks)
}

func (controller *AdminPalletRackController) parseInventoryID(c *gin.Context) (*uint, error) {
	inventoryIDStr := c.Query("inventory_id")
	if inventoryIDStr == "" {
		return nil, nil
	}

	id, err := strconv.ParseUint(inventoryIDStr, 10, 32)
	if err != nil {
		return nil, err
	}

	idUint := uint(id)
	return &idUint, nil
}

func (controller *AdminPalletRackController) parsePage(c *gin.Context) int {
	pageStr := c.Query("page")
	if pageStr == "" {
		return 1
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		return 1
	}

	return page
}

func (controller *AdminPalletRackController) parseLimit(c *gin.Context) int {
	limitStr := c.Query("limit")
	if limitStr == "" {
		return 10
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		return 10
	}

	return limit
}

// @Summary Admin List Pallet Racks
// @Tags Admin Pallet Racks
// @Accept json
// @Produce json
// @Param inventory_id query int false "Inventory ID to filter"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} apiContracts.PaginatedRacksResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/admin/racks [get]
func (controller *AdminPalletRackController) ListRacks(c *gin.Context) {
	inventoryID, err := controller.parseInventoryID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inventory_id"})
		return
	}

	page := controller.parsePage(c)
	limit := controller.parseLimit(c)

	racks, err := controller.palletRackService.ListRacks(inventoryID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, racks)
}
