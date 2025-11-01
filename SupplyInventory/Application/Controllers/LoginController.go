package controllers

import (
	"net/http"

	requestsHelper "go_inventory/Helpers/RequestsHelper"
	requests "go_inventory/SupplyInventory/Application/Requests"
	services "go_inventory/SupplyInventory/Application/Services"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	jwtService  services.JWTService
	userService services.UserService
}

func NewLoginController(jwtService services.JWTService, userService services.UserService) *LoginController {
	return &LoginController{jwtService: jwtService, userService: userService}
}

func (controller *LoginController) RegisterLogin(group *gin.RouterGroup) {
	group.POST("/login", controller.Login)
}

func (controller *LoginController) Login(context *gin.Context) {
	var req requests.LoginRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": requestsHelper.FormatValidationErrors(err)})
		return
	}

	user, err := controller.userService.Login(req.Email, req.Password)
	if err != nil {
		context.JSON(err.ErrorCode(), err.Error())
		return
	}

	token, appErr := controller.jwtService.GenerateToken(user.ID, user.Email)
	if appErr != nil {
		context.JSON(appErr.ErrorCode(), appErr.Error())
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": token})
}
