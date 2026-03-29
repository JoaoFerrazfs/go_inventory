//go:build unit

package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	middlewares "go_inventory/SupplyInventory/Application/Middlewares"
	"go_inventory/SupplyInventory/Domain/constants"
)

// TestRBACMiddleware_RequireAdmin_Success verifica que admin consegue acessar endpoint protegido
func TestRBACMiddleware_RequireAdmin_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Middleware de teste para simular autenticação com role
	router.Use(func(c *gin.Context) {
		c.Set("role", constants.AdminRole)
		c.Next()
	})

	// Aplicar middleware RBAC
	router.Use(middlewares.RequireAdmin())

	router.GET("/admin-only", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
	})

	// Action
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin-only", nil)
	router.ServeHTTP(w, req)

	// Verify
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestRBACMiddleware_RequireAdmin_Forbidden verifica que user não consegue acessar endpoint protegido
func TestRBACMiddleware_RequireAdmin_Forbidden(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Middleware de teste para simular autenticação com role user
	router.Use(func(c *gin.Context) {
		c.Set("role", constants.RoleUser)
		c.Next()
	})

	// Aplicar middleware RBAC
	router.Use(middlewares.RequireAdmin())

	router.GET("/admin-only", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
	})

	// Action
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin-only", nil)
	router.ServeHTTP(w, req)

	// Verify
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// TestRBACMiddleware_RequireRole_MultipleRoles verifica múltiplas roles permitidas
func TestRBACMiddleware_RequireRole_MultipleRoles(t *testing.T) {
	tests := []struct {
		name       string
		role       constants.UserRole
		shouldPass bool
	}{
		{
			name:       "Admin consegue acessar rota que permite admin ou user",
			role:       constants.AdminRole,
			shouldPass: true,
		},
		{
			name:       "User consegue acessar rota que permite admin ou user",
			role:       constants.RoleUser,
			shouldPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			router := gin.New()

			// Middleware de teste para simular autenticação
			router.Use(func(c *gin.Context) {
				c.Set("role", tt.role)
				c.Next()
			})

			// Aplicar middleware RBAC para múltiplas roles
			router.Use(middlewares.RequireRole(constants.AdminRole, constants.RoleUser))

			router.GET("/protected", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "access granted"})
			})

			// Action
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/protected", nil)
			router.ServeHTTP(w, req)

			// Verify
			if tt.shouldPass {
				assert.Equal(t, http.StatusOK, w.Code)
			} else {
				assert.Equal(t, http.StatusForbidden, w.Code)
			}
		})
	}
}

// TestRBACMiddleware_NoRoleSet verifica comportamento quando role não está setada
func TestRBACMiddleware_NoRoleSet(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// NÃO setar role - simular requisição sem autenticação

	// Aplicar middleware RBAC
	router.Use(middlewares.RequireAdmin())

	router.GET("/admin-only", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
	})

	// Action
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin-only", nil)
	router.ServeHTTP(w, req)

	// Verify - deve retornar 403 pois role não foi setada
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// TestRBACMiddleware_InvalidRoleType verifica comportamento com tipo inválido
func TestRBACMiddleware_InvalidRoleType(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Middleware que seta role com tipo inválido
	router.Use(func(c *gin.Context) {
		c.Set("role", "invalid_type_not_UserRole") // String simples, não UserRole
		c.Next()
	})

	// Aplicar middleware RBAC
	router.Use(middlewares.RequireAdmin())

	router.GET("/admin-only", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
	})

	// Action
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin-only", nil)
	router.ServeHTTP(w, req)

	// Verify - deve retornar 403 pois tipo é inválido
	assert.Equal(t, http.StatusForbidden, w.Code)
}
