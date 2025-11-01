package controllers

import (
	"net/http"

	requests "go_inventory/SupplyInventory/Application/Requests"
	services "go_inventory/SupplyInventory/Application/Services"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	jwtService services.JWTService
}

func NewLoginController(jwtService services.JWTService) *LoginController {
	return &LoginController{jwtService: jwtService}
}

func (controller *LoginController) RegisterLogin(group *gin.RouterGroup) {
	group.POST("/login", controller.Login)
}

func (controller *LoginController) Login(context *gin.Context) {
	var req requests.LoginRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.Username != "admin" || req.Password != "1234" {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, appErr := controller.jwtService.GenerateToken(1234, req.Username)
	if appErr != nil {
		context.JSON(appErr.ErrorCode(), appErr.Error())
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": token})
}
