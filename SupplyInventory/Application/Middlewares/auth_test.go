package middlewares_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	errors "go_inventory/Helpers/Errors"
	middlewares "go_inventory/SupplyInventory/Application/Middlewares"
	"go_inventory/SupplyInventory/Domain/constants"

	"github.com/gin-gonic/gin"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

type stubJWTService struct {
	token  *jwtv5.Token
	appErr *errors.AppError
}

func (s *stubJWTService) GenerateToken(userID uint, username string, role constants.UserRole) (string, *errors.AppError) {
	return "", nil
}

func (s *stubJWTService) GenerateRefreshToken(userID uint, username string, role constants.UserRole) (string, *errors.AppError) {
	return "", nil
}

func (s *stubJWTService) ValidateToken(tokenString string) (*jwtv5.Token, *errors.AppError) {
	return s.token, s.appErr
}

func (s *stubJWTService) RefreshToken(token string) (string, *errors.AppError) {
	return "", nil
}

func performRequest(handler http.HandlerFunc, method, path string, header map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	handler(w, req)
	return w
}

// newRequestWithBearer returns a request with an Authorization header using the
// `Bearer <token>` format. Tests use a stubbed JWT service so the string value
// itself is irrelevant; this helper makes the header source explicit.
func newRequestWithBearer(method, path, token string) *http.Request {
	req := httptest.NewRequest(method, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

// newRequestWithAuthHeader returns a request with a custom Authorization header value.
// Useful for testing invalid formats like `Token abc`.
func newRequestWithAuthHeader(method, path, headerValue string) *http.Request {
	req := httptest.NewRequest(method, path, nil)
	if headerValue != "" {
		req.Header.Set("Authorization", headerValue)
	}
	return req
}

func TestAuth_MissingHeader(t *testing.T) {
	// Set
	gin.SetMode(gin.TestMode)
	m := &middlewares.AuthMiddleware{JWTService: &stubJWTService{}}
	r := gin.New()
	r.Use(m.Handler())
	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// Actions
	// Actions
	w := httptest.NewRecorder()
	req := newRequestWithBearer(http.MethodGet, "/", "")
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, "missing token", body["error"])
}

func TestAuth_InvalidFormat(t *testing.T) {
	// Set
	gin.SetMode(gin.TestMode)
	m := &middlewares.AuthMiddleware{JWTService: &stubJWTService{}}
	r := gin.New()
	r.Use(m.Handler())
	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// Actions
	// Actions
	w := httptest.NewRecorder()
	req := newRequestWithAuthHeader(http.MethodGet, "/", "Token abc")
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, "invalid token format", body["error"])
}

func TestAuth_ValidateError(t *testing.T) {
	// Set
	gin.SetMode(gin.TestMode)
	stub := &stubJWTService{token: nil, appErr: errors.NewAppError("some", 401)}
	m := &middlewares.AuthMiddleware{JWTService: stub}
	r := gin.New()
	r.Use(m.Handler())
	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// Actions
	// Actions
	w := httptest.NewRecorder()
	req := newRequestWithBearer(http.MethodGet, "/", "token")
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, "invalid token", body["error"])
}

func TestAuth_InvalidTokenType(t *testing.T) {
	// Set
	gin.SetMode(gin.TestMode)
	token := &jwtv5.Token{Valid: true, Claims: jwtv5.MapClaims{"userID": float64(1), "username": "u", "tokenType": "refresh"}}
	stub := &stubJWTService{token: token, appErr: nil}
	m := &middlewares.AuthMiddleware{JWTService: stub}
	r := gin.New()
	r.Use(m.Handler())
	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// Actions
	// Actions
	w := httptest.NewRecorder()
	req := newRequestWithBearer(http.MethodGet, "/", "token")
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	// body is AppError JSON
	var body map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, "invalid token type", body["message"])
}

func TestAuth_Success(t *testing.T) {
	// Set
	gin.SetMode(gin.TestMode)
	token := &jwtv5.Token{Valid: true, Claims: jwtv5.MapClaims{"userID": float64(42), "username": "alice", "tokenType": "access", "role": string(constants.RoleUser)}}
	stub := &stubJWTService{token: token, appErr: nil}
	m := &middlewares.AuthMiddleware{JWTService: stub}
	r := gin.New()
	r.Use(m.Handler())
	r.GET("/", func(c *gin.Context) {
		uid, _ := c.Get("userID")
		uname, _ := c.Get("username")
		c.JSON(200, gin.H{"userID": uid, "username": uname})
	})

	// Actions
	// Actions
	w := httptest.NewRecorder()
	req := newRequestWithBearer(http.MethodGet, "/", "token")
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	// Gin serializes numbers as float64 in JSON
	assert.Equal(t, float64(42), body["userID"])
	assert.Equal(t, "alice", body["username"])
}
