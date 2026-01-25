package services_test

import (
	"testing"

	errors "go_inventory/Helpers/Errors"
	productService "go_inventory/SupplyInventory/Application/Services/Product"
	entities "go_inventory/SupplyInventory/Domain/Entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockProductRepo struct{ mock.Mock }

func (m *mockProductRepo) Create(ean string, name string) (*entities.ProductEntity, *errors.AppError) {
	args := m.Called(ean, name)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*entities.ProductEntity), nil
}

func (m *mockProductRepo) FindByEAN(ean string) (*entities.ProductEntity, *errors.AppError) {
	args := m.Called(ean)
	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil
		}
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*entities.ProductEntity), nil
}

func (m *mockProductRepo) Delete(ean string) (bool, *errors.AppError) {
	args := m.Called(ean)
	if args.Get(1) != nil {
		return args.Bool(0), args.Get(1).(*errors.AppError)
	}
	return args.Bool(0), nil
}

func TestCreateProduct_Success(t *testing.T) {
	// Set
	repo := &mockProductRepo{}
	svc := productService.NewProductService(repo)

	// Expectations
	repo.On("FindByEAN", "1234567890123").Return(nil, nil)
	created := &entities.ProductEntity{EAN: "1234567890123", Name: "Test Product"}
	repo.On("Create", "1234567890123", "Test Product").Return(created, nil)

	// Actions
	res, err := svc.CreateProduct("1234567890123", "Test Product")

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, created, res)
	repo.AssertExpectations(t)
}

func TestCreateProduct_DuplicateEAN(t *testing.T) {
	// Set
	repo := &mockProductRepo{}
	svc := productService.NewProductService(repo)

	// Expectations
	existing := &entities.ProductEntity{EAN: "1234567890123", Name: "Existing Product"}
	repo.On("FindByEAN", "1234567890123").Return(existing, nil)

	// Actions
	res, err := svc.CreateProduct("1234567890123", "New Product")

	// Assertions
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, 400, err.Code)
	assert.Equal(t, "Product with given EAN already exists", err.Message)
	repo.AssertExpectations(t)
}

func TestCreateProduct_RepoCreateError(t *testing.T) {
	// Set
	repo := &mockProductRepo{}
	svc := productService.NewProductService(repo)

	// Expectations
	repo.On("FindByEAN", "1234567890123").Return(nil, nil)
	appErr := errors.NewAppError("Database error", 500)
	repo.On("Create", "1234567890123", "Test Product").Return(nil, appErr)

	// Actions
	res, err := svc.CreateProduct("1234567890123", "Test Product")

	// Assertions
	assert.Nil(t, res)
	assert.Equal(t, appErr, err)
	repo.AssertExpectations(t)
}

func TestGetProductByEAN_Success(t *testing.T) {
	// Set
	repo := &mockProductRepo{}
	svc := productService.NewProductService(repo)

	// Expectations
	product := &entities.ProductEntity{EAN: "1234567890123", Name: "Test Product"}
	repo.On("FindByEAN", "1234567890123").Return(product, nil)

	// Actions
	res, err := svc.GetProductByEAN("1234567890123")

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, product, res)
	repo.AssertExpectations(t)
}

func TestGetProductByEAN_NotFound(t *testing.T) {
	// Set
	repo := &mockProductRepo{}
	svc := productService.NewProductService(repo)

	// Expectations
	appErr := errors.NewAppError("Product not found", 404)
	repo.On("FindByEAN", "1234567890123").Return(nil, appErr)

	// Actions
	res, err := svc.GetProductByEAN("1234567890123")

	// Assertions
	assert.Nil(t, res)
	assert.Equal(t, appErr, err)
	repo.AssertExpectations(t)
}

func TestDeleteProduct_Success(t *testing.T) {
	// Set
	repo := &mockProductRepo{}
	svc := productService.NewProductService(repo)

	// Expectations
	repo.On("Delete", "1234567890123").Return(true, nil)

	// Actions
	ok, err := svc.DeleteProduct("1234567890123")

	// Assertions
	assert.Nil(t, err)
	assert.True(t, ok)
	repo.AssertExpectations(t)
}

func TestDeleteProduct_Error(t *testing.T) {
	// Set
	repo := &mockProductRepo{}
	svc := productService.NewProductService(repo)

	// Expectations
	appErr := errors.NewAppError("Delete failed", 500)
	repo.On("Delete", "1234567890123").Return(false, appErr)

	// Actions
	ok, err := svc.DeleteProduct("1234567890123")

	// Assertions
	assert.False(t, ok)
	assert.Equal(t, appErr, err)
	repo.AssertExpectations(t)
}
