package controllers

import (
	"net/http"

	requestsHelper "go_inventory/Helpers/RequestsHelper"
	requests "go_inventory/SupplyInventory/Application/Requests"
	responses "go_inventory/SupplyInventory/Application/Responses"
	jwt "go_inventory/SupplyInventory/Application/Services/Jwt"
	user "go_inventory/SupplyInventory/Application/Services/User"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	jwtService  jwt.JWTService
	userService user.UserService
}

func NewAuthController(jwtService jwt.JWTService, userService user.UserService) *AuthController {
	return &AuthController{jwtService: jwtService, userService: userService}
}

func (controller *AuthController) RegisterLogin(group *gin.RouterGroup) {
	group.POST("/login", controller.Login)
	group.POST("/refreshToken", controller.RefreshToken)
}

// @Summary Login
// @Tags Authentication
// @Accept json
// @Produce json
// @Param userData body requests.LoginRequest true "User Data"
// @Success 200 {object} responses.AuthResponse
// @Failure 422 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/auth/login [post]
func (controller *AuthController) Login(context *gin.Context) {
	var req requests.LoginRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": requestsHelper.FormatValidationErrors(err)})
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

	refreshToken, appErr := controller.jwtService.GenerateRefreshToken(user.ID, user.Email)
	if appErr != nil {
		context.JSON(appErr.ErrorCode(), appErr.Error())
		return
	}

	response := responses.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}

	context.JSON(http.StatusOK, response)
}

// @Summary Refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param refreshToken body requests.RefreshTokenRequest true "Refresh Token"
// @Success 200 {object} responses.RefreshResponse
// @Failure 401 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/auth/refreshToken [post]
func (controller *AuthController) RefreshToken(context *gin.Context) {
	var req requests.RefreshTokenRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": requestsHelper.FormatValidationErrors(err)})
		return
	}

	token, appErr := controller.jwtService.RefreshToken(req.RefreshToken)
	if appErr != nil {
		context.JSON(appErr.ErrorCode(), appErr.Error())
		return
	}

	response := responses.RefreshResponse{
		Token: token,
	}

	context.JSON(http.StatusOK, response)
}
