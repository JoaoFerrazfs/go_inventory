package middlewares_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	errors "go_inventory/Helpers/Errors"
	middlewares "go_inventory/SupplyInventory/Application/Middlewares"
	entities "go_inventory/SupplyInventory/Domain/Entities"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type stubInventoryRepo struct {
	exists bool
}

func (s *stubInventoryRepo) Exists(id uint) (bool, *errors.AppError) {
	if s.exists {
		return true, nil
	}
	return false, errors.NewAppError("Inventory not found", 404)
}

func (s *stubInventoryRepo) Where(conditions ...map[string]any) ([]entities.InventoryEntity, *errors.AppError) {
	return []entities.InventoryEntity{}, nil
}

func (s *stubInventoryRepo) FindById(id uint) (*entities.InventoryEntity, *errors.AppError) {
	return &entities.InventoryEntity{ID: id}, nil
}

func (s *stubInventoryRepo) Create(inventory *entities.InventoryEntity) *errors.AppError {
	return nil
}

func (s *stubInventoryRepo) Update(inventory *entities.InventoryEntity) *errors.AppError {
	return nil
}

func TestInventory_MissingHeader(t *testing.T) {
	// Set
	gin.SetMode(gin.TestMode)
	repo := &stubInventoryRepo{exists: false}
	m := middlewares.NewInventoryMiddleware(repo)
	r := gin.New()
	r.Use(m.Handler())
	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// Actions
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 422, w.Code)
	var body map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, float64(422), body["code"])
}

func TestInventory_InvalidHeader(t *testing.T) {
	// Set
	gin.SetMode(gin.TestMode)
	repo := &stubInventoryRepo{exists: false}
	m := middlewares.NewInventoryMiddleware(repo)
	r := gin.New()
	r.Use(m.Handler())
	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// Actions
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Inventory-ID", "abc")
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 422, w.Code)
	var body map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, float64(422), body["code"])
}

func TestInventory_NotFound(t *testing.T) {
	// Set
	gin.SetMode(gin.TestMode)
	repo := &stubInventoryRepo{exists: false}
	m := middlewares.NewInventoryMiddleware(repo)
	r := gin.New()
	r.Use(m.Handler())
	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// Actions
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Inventory-ID", "123")
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 404, w.Code)
	var body map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, float64(404), body["code"])
}

func TestInventory_Success(t *testing.T) {
	// Set
	gin.SetMode(gin.TestMode)
	repo := &stubInventoryRepo{exists: true}
	m := middlewares.NewInventoryMiddleware(repo)
	r := gin.New()
	r.Use(m.Handler())
	r.GET("/", func(c *gin.Context) { id := middlewares.GetInventoryID(c); c.JSON(200, gin.H{"id": id}) })

	// Actions
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Inventory-ID", "123")
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 200, w.Code)
}
