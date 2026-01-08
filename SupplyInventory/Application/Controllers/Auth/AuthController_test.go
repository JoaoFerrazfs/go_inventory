package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	errors "go_inventory/Helpers/Errors"

	auth "go_inventory/SupplyInventory/Application/Controllers/Auth"
	entities "go_inventory/SupplyInventory/Domain/Entities"
)

type mockJWTService struct {
	mock.Mock
}

func (m *mockJWTService) GenerateToken(userID uint, email string) (string, *errors.AppError) {
	args := m.Called(userID, email)
	var appErr *errors.AppError
	if args.Get(1) != nil {
		appErr = args.Get(1).(*errors.AppError)
	}
	return args.String(0), appErr
}

func (m *mockJWTService) GenerateRefreshToken(userID uint, email string) (string, *errors.AppError) {
	args := m.Called(userID, email)
	var appErr *errors.AppError
	if args.Get(1) != nil {
		appErr = args.Get(1).(*errors.AppError)
	}
	return args.String(0), appErr
}

func (m *mockJWTService) ValidateToken(tokenString string) (*jwt.Token, *errors.AppError) {
	return nil, nil
}

func (m *mockJWTService) RefreshToken(refreshToken string) (string, *errors.AppError) {
	args := m.Called(refreshToken)
	var appErr *errors.AppError
	if args.Get(1) != nil {
		appErr = args.Get(1).(*errors.AppError)
	}
	return args.String(0), appErr
}

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) CreateUser(name, email, password string) (*entities.UserEntity, *errors.AppError) {
	return nil, nil
}

func (m *mockUserService) Login(email, password string) (*entities.UserEntity, *errors.AppError) {
	args := m.Called(email, password)
	var appErr *errors.AppError
	if args.Get(1) != nil {
		appErr = args.Get(1).(*errors.AppError)
	}
	if args.Get(0) == nil {
		return nil, appErr
	}
	return args.Get(0).(*entities.UserEntity), appErr
}

func (m *mockUserService) GetUserByID(id uint) (*entities.UserEntity, *errors.AppError) {
	args := m.Called(id)
	var appErr *errors.AppError
	if args.Get(1) != nil {
		appErr = args.Get(1).(*errors.AppError)
	}
	if args.Get(0) == nil {
		return nil, appErr
	}
	return args.Get(0).(*entities.UserEntity), appErr
}

func TestLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set
	jwtMock := new(mockJWTService)
	userMock := new(mockUserService)
	controller := auth.NewAuthController(jwtMock, userMock)

	user := &entities.UserEntity{ID: 1, Email: "test@example.com"}

	// Expectations
	userMock.On("Login", "test@example.com", "password").Return(user, nil)

	jwtMock.On("GenerateToken", user.ID, user.Email).Return("token", nil)

	jwtMock.On("GenerateRefreshToken", user.ID, user.Email).Return("refresh", nil)

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]string{"email": "test@example.com", "password": "password"})

	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.Login(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	// Assertions - response body
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "token", resp["token"])
	assert.Equal(t, "refresh", resp["refreshToken"])
}

func TestLogin_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set
	jwtMock := new(mockJWTService)
	userMock := new(mockUserService)
	controller := auth.NewAuthController(jwtMock, userMock)

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := []byte(`{"email":123}`)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.Login(c)

	// Assertions
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestLogin_UserServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set
	jwtMock := new(mockJWTService)
	userMock := new(mockUserService)
	controller := auth.NewAuthController(jwtMock, userMock)

	// Expectations
	userMock.On("Login", "fail@example.com", "bad").Return(nil, errors.NewAppError("invalid", 401))

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body, _ := json.Marshal(map[string]string{"email": "fail@example.com", "password": "bad"})
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.Login(c)

	// Assertions
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestRefreshToken_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set
	jwtMock := new(mockJWTService)
	userMock := new(mockUserService)
	controller := auth.NewAuthController(jwtMock, userMock)

	// Expectations
	jwtMock.On("RefreshToken", "refresh_token").Return("new_token", nil)

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body, _ := json.Marshal(map[string]string{"refreshToken": "refresh_token"})
	c.Request, _ = http.NewRequest("POST", "/refreshToken", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.RefreshToken(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRefreshToken_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set
	jwtMock := new(mockJWTService)
	userMock := new(mockUserService)
	controller := auth.NewAuthController(jwtMock, userMock)

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := []byte(`{"refreshToken":123}`)
	c.Request, _ = http.NewRequest("POST", "/refreshToken", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.RefreshToken(c)

	// Assertions
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestRefreshToken_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set
	jwtMock := new(mockJWTService)
	userMock := new(mockUserService)
	controller := auth.NewAuthController(jwtMock, userMock)

	// Expectations
	jwtMock.On("RefreshToken", "bad_token").Return("", errors.NewAppError("invalid", 401))

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body, _ := json.Marshal(map[string]string{"refreshToken": "bad_token"})
	c.Request, _ = http.NewRequest("POST", "/refreshToken", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.RefreshToken(c)

	// Assertions
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestLogin_EmptyTokenAndRefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set
	jwtMock := new(mockJWTService)
	userMock := new(mockUserService)
	controller := auth.NewAuthController(jwtMock, userMock)

	user := &entities.UserEntity{ID: 4, Email: "empty@example.com"}

	// Expectations
	userMock.On("Login", "empty@example.com", "password").Return(user, nil)
	jwtMock.On("GenerateToken", user.ID, user.Email).Return("", nil)
	jwtMock.On("GenerateRefreshToken", user.ID, user.Email).Return("", nil)

	// Actions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body, _ := json.Marshal(map[string]string{"email": "empty@example.com", "password": "password"})
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.Login(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "", resp["token"])
	assert.Equal(t, "", resp["refreshToken"])
}
