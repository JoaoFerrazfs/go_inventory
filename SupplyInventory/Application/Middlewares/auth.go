package middlewares

import (
	"net/http"
	"strings"

	services "go_inventory/SupplyInventory/Application/Services"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	JWTService services.JWTService
}

func NewAuthMiddleware(jwtService services.JWTService) *AuthMiddleware {
	return &AuthMiddleware{JWTService: jwtService}
}

func (m *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}

		valid, appErr := m.JWTService.ValidateToken(tokenString)
		if appErr != nil || !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Next()
	}
}
