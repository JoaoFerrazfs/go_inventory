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

	errors "go_inventory/Helpers/Errors"
	user "go_inventory/SupplyInventory/Application/Controllers/User"
	entities "go_inventory/SupplyInventory/Domain/Entities"
)

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) CreateUser(name, email, password string) (*entities.UserEntity, *errors.AppError) {
	args := m.Called(name, email, password)
	var appErr *errors.AppError
	if args.Get(1) != nil {
		appErr = args.Get(1).(*errors.AppError)
	}
	if args.Get(0) == nil {
		return nil, appErr
	}
	return args.Get(0).(*entities.UserEntity), appErr
}

func (m *mockUserService) Login(email, password string) (*entities.UserEntity, *errors.AppError) {
	return nil, nil
}

func TestCreateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set
	userMock := new(mockUserService)
	controller := user.NewUserController(userMock)
	user := &entities.UserEntity{ID: 1, Email: "test@example.com", Name: "Test User"}

	// Expectations
	userMock.On("CreateUser", "Test User", "test@example.com", "password").Return(user, nil)

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body, _ := json.Marshal(map[string]string{"name": "Test User", "email": "test@example.com", "password": "password"})
	c.Request, _ = http.NewRequest("POST", "/create", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.Create(c)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateUser_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set
	userMock := new(mockUserService)
	controller := user.NewUserController(userMock)

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := []byte(`{"name":123}`)
	c.Request, _ = http.NewRequest("POST", "/create", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.Create(c)

	// Assertions
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestCreateUser_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set
	userMock := new(mockUserService)
	controller := user.NewUserController(userMock)

	// Expectations
	userMock.On("CreateUser", "Fail User", "fail@example.com", "bad").Return(nil, errors.NewAppError("fail", 500))

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body, _ := json.Marshal(map[string]string{"name": "Fail User", "email": "fail@example.com", "password": "bad"})
	c.Request, _ = http.NewRequest("POST", "/create", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.Create(c)

	// Assertions
	assert.NotEqual(t, http.StatusCreated, w.Code)
}
