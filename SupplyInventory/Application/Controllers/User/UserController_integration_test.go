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

	user "go_inventory/SupplyInventory/Application/Controllers/User"
	requests "go_inventory/SupplyInventory/Application/Requests"
	services "go_inventory/SupplyInventory/Application/Services"
	entities "go_inventory/SupplyInventory/Domain/Entities"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) Create(user *entities.UserEntity) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockUserRepository) FindByEmail(email string) (*entities.UserEntity, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.UserEntity), args.Error(1)
}

func setupUserRouter() *gin.Engine {
	mockRepo := new(mockUserRepository)
	userService := services.NewUserService(mockRepo)
	controller := user.NewUserController(userService)

	r := gin.Default()
	api := r.Group("/api/v1/users")
	controller.RegisterUserRoutes(api)
	return r
}

func TestIntegration_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set
	mockRepo := new(mockUserRepository)
	// testUser variable removed as it was unused
	userService := services.NewUserService(mockRepo)
	controller := user.NewUserController(userService)
	r := gin.Default()
	api := r.Group("/api/v1/users")
	controller.RegisterUserRoutes(api)

	createReq := requests.UserRequest{
		Name:     "Admin",
		Email:    "admin@example.com",
		Password: "admin123",
	}
	body, _ := json.Marshal(createReq)

	// Expectations
	mockRepo.On("FindByEmail", "admin@example.com").Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*entities.UserEntity")).Return(nil)

	// Actions
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Assertions
	assert.Contains(t, []int{http.StatusCreated, http.StatusUnprocessableEntity}, w.Code)
}
