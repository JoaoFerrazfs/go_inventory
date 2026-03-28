package middlewares

import (
	"net/http"

	"go_inventory/SupplyInventory/Domain/constants"

	"github.com/gin-gonic/gin"
)

type RBACMiddleware struct{}

func NewRBACMiddleware() *RBACMiddleware {
	return &RBACMiddleware{}
}

// RequireRole verifica se o usuário tem uma das roles especificadas.
// Chamadas são feitas após o middleware de autenticação
func (m *RBACMiddleware) RequireRole(allowedRoles ...constants.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "user role not found in token",
			})
			return
		}

		role, ok := userRole.(constants.UserRole)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "invalid role format",
			})
			return
		}

		// Verifica se a role do usuário está na lista de roles permitidas
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "insufficient permissions",
		})
	}
}

// RequireAdmin é um atalho para RequireRole(AdminRole)
func (m *RBACMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole(constants.AdminRole)
}

// RequireAny verifica se o usuário tem ao menos uma das roles
func (m *RBACMiddleware) RequireAny(roles ...constants.UserRole) gin.HandlerFunc {
	return m.RequireRole(roles...)
}

// Convenience functions for direct use
func RequireRole(allowedRoles ...constants.UserRole) gin.HandlerFunc {
	return NewRBACMiddleware().RequireRole(allowedRoles...)
}

func RequireAdmin() gin.HandlerFunc {
	return NewRBACMiddleware().RequireAdmin()
}

func RequireAny(roles ...constants.UserRole) gin.HandlerFunc {
	return NewRBACMiddleware().RequireAny(roles...)
}
