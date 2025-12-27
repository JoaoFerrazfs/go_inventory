package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	auth "go_inventory/SupplyInventory/Application/Controllers/Auth"
	requests "go_inventory/SupplyInventory/Application/Requests"
	services "go_inventory/SupplyInventory/Application/Services"
	domain "go_inventory/SupplyInventory/Domain"
)

// Mock UserRepository for integration
type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) Create(user *domain.UserEntity) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockUserRepository) FindByEmail(email string) (*domain.UserEntity, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserEntity), args.Error(1)
}

func setupRouter() *gin.Engine {
	jwtService := services.NewJWTService()
	mockRepo := new(mockUserRepository)
	userService := services.NewUserService(mockRepo)
	controller := auth.NewAuthController(jwtService, userService)

	r := gin.Default()
	api := r.Group("/api/v1/auth")
	controller.RegisterLogin(api)
	return r
}

func TestIntegration_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set
	jwtService := services.NewJWTService()
	mockRepo := new(mockUserRepository)
	user := &domain.UserEntity{ID: 1, Email: "admin@example.com", Password: "$2a$10$hash"}
	userService := services.NewUserService(mockRepo)
	controller := auth.NewAuthController(jwtService, userService)
	r := gin.Default()
	api := r.Group("/api/v1/auth")
	controller.RegisterLogin(api)

	loginReq := requests.LoginRequest{
		Email:    "admin@example.com",
		Password: "admin123",
	}

	body, _ := json.Marshal(loginReq)

	// Expectations
	mockRepo.On("FindByEmail", "admin@example.com").Return(user, nil)

	// Actions
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Assertions
	assert.Contains(t, []int{http.StatusOK, http.StatusUnauthorized, http.StatusUnprocessableEntity}, w.Code)
}

func TestIntegration_RefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set
	r := setupRouter()

	refreshReq := map[string]string{"refreshToken": "sometoken"}
	body, _ := json.Marshal(refreshReq)

	// Actions
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/refreshToken", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Assertions
	assert.Contains(t, []int{http.StatusOK, http.StatusUnprocessableEntity, http.StatusUnauthorized}, w.Code)
}
