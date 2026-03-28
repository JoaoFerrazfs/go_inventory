package middlewares

import (
	"net/http"
	"strings"

	errors "go_inventory/Helpers/Errors"
	jwt "go_inventory/SupplyInventory/Application/Services/Jwt"
	"go_inventory/SupplyInventory/Domain/constants"

	"github.com/gin-gonic/gin"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	JWTService jwt.JWTService
}

func NewAuthMiddleware(jwtService jwt.JWTService) *AuthMiddleware {
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

		token, appErr := m.JWTService.ValidateToken(tokenString)
		if appErr != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwtv5.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.NewAppError("invalid token", 401))
			return
		}

		tokenType, ok := claims["tokenType"].(string)
		if !ok || tokenType != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.NewAppError("invalid token type", 403))
			return
		}

		c.Set("userID", uint(claims["userID"].(float64)))
		c.Set("username", claims["username"].(string))
		c.Set("role", constants.UserRole(claims["role"].(string)))

		c.Next()
	}
}
