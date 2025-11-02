package controllers

import (
	"net/http"

	requestsHelper "go_inventory/Helpers/RequestsHelper"
	requests "go_inventory/SupplyInventory/Application/Requests"
	services "go_inventory/SupplyInventory/Application/Services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (controller *UserController) RegisterUserRoutes(group *gin.RouterGroup) {
	group.POST("/create", controller.create)
}

// @Summary Create User
// @Tags Users
// @Accept json
// @Produce json
// @Param user body requests.UserRequest true "User"
// @Success 201 {object} domain.UserEntity
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users/create [post]
func (controller *UserController) create(c *gin.Context) {
	var req requests.UserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": requestsHelper.FormatValidationErrors(err)})
		return
	}

	user, appErr := controller.userService.CreateUser(req.Name, req.Email, req.Password)
	if appErr != nil {
		c.JSON(appErr.Code, appErr.Message)
		return
	}

	c.JSON(http.StatusCreated, user)
}
