package services

import (
	"testing"

	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPalletRepoPP struct{ mock.Mock }

func (m *mockPalletRepoPP) GetSupplyById(id uint) (*entities.PalletEntity, *errors.AppError) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*entities.PalletEntity), nil
}

// Additional repository methods to satisfy interface
func (m *mockPalletRepoPP) AddProductsToPallet(product entities.PalletizedProductEntity) (*entities.PalletEntity, *errors.AppError) {
	args := m.Called(product)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*entities.PalletEntity), nil
}
func (m *mockPalletRepoPP) GetAllPallets(palletRackId *uint, productId *uint) ([]entities.PalletEntity, *errors.AppError) {
	return nil, nil
}
func (m *mockPalletRepoPP) AddSupply(name string, rackId uint) (*entities.PalletEntity, *errors.AppError) {
	return nil, nil
}
func (m *mockPalletRepoPP) UpdateSupply(p *entities.PalletEntity) (*entities.PalletEntity, *errors.AppError) {
	return nil, nil
}
func (m *mockPalletRepoPP) DeletePalletById(id uint) (bool, *errors.AppError) { return false, nil }
func (m *mockPalletRepoPP) Create(pallet *entities.PalletEntity) error        { return nil }
func (m *mockPalletRepoPP) FindByID(id uint) (*entities.PalletEntity, error)  { return nil, nil }
func (m *mockPalletRepoPP) List() ([]*entities.PalletEntity, error)           { return nil, nil }
func (m *mockPalletRepoPP) DeleteByID(id uint) error                          { return nil }
func (m *mockPalletRepoPP) Update(pallet *entities.PalletEntity) error        { return nil }

type mockPalletizedProductRepo struct{ mock.Mock }

func (m *mockPalletizedProductRepo) AddProductsToPallet(product entities.PalletizedProductEntity) (bool, *errors.AppError) {
	args := m.Called(product)
	if args.Get(1) != nil {
		return args.Bool(0), args.Get(1).(*errors.AppError)
	}
	return args.Bool(0), nil
}
func (m *mockPalletizedProductRepo) DeleteProductsFromPallet(palletId uint, productsEan int) (bool, *errors.AppError) {
	args := m.Called(palletId, productsEan)
	if args.Get(1) != nil {
		return args.Bool(0), args.Get(1).(*errors.AppError)
	}
	return args.Bool(0), nil
}

func TestAddProductsToPallet_Success(t *testing.T) {
	// Set
	palletRepo := &mockPalletRepoPP{}
	repoProd := &mockPalletizedProductRepo{}
	svc := NewPalletizedProductService(palletRepo, repoProd)

	// Expectations
	repoProd.On("AddProductsToPallet", mock.Anything).Return(true, nil)
	expected := &entities.PalletEntity{ID: 1}
	palletRepo.On("GetSupplyById", uint(1)).Return(expected, nil)

	// Actions
	res, err := svc.AddProductsToPallet(1, 12345, 2)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, expected, res)
	repoProd.AssertExpectations(t)
	palletRepo.AssertExpectations(t)
}

func TestAddProductsToPallet_AddError(t *testing.T) {
	// Set
	palletRepo := &mockPalletRepoPP{}
	repoProd := &mockPalletizedProductRepo{}
	svc := NewPalletizedProductService(palletRepo, repoProd)

	// Expectations
	appErr := errors.NewAppError("add error", 400)
	repoProd.On("AddProductsToPallet", mock.Anything).Return(false, appErr)

	// Actions
	res, err := svc.AddProductsToPallet(2, 11111, 1)

	// Assertions
	assert.Nil(t, res)
	assert.NotNil(t, err)
	repoProd.AssertExpectations(t)
}

func TestDeleteProductsFromPallet_Success(t *testing.T) {
	// Set
	palletRepo := &mockPalletRepoPP{}
	repoProd := &mockPalletizedProductRepo{}
	svc := NewPalletizedProductService(palletRepo, repoProd)

	// Expectations
	repoProd.On("DeleteProductsFromPallet", uint(1), 12345).Return(true, nil)

	// Actions
	deleted, err := svc.DeleteProductsFromPallet(1, 12345)

	// Assertions
	assert.Nil(t, err)
	assert.True(t, deleted)
	repoProd.AssertExpectations(t)
}

func TestDeleteProductsFromPallet_RepoError(t *testing.T) {
	// Set
	palletRepo := &mockPalletRepoPP{}
	repoProd := &mockPalletizedProductRepo{}
	svc := NewPalletizedProductService(palletRepo, repoProd)

	// Expectations
	appErr := errors.NewAppError("delete error", 400)
	repoProd.On("DeleteProductsFromPallet", uint(2), 11111).Return(false, appErr)

	// Actions
	deleted, err := svc.DeleteProductsFromPallet(2, 11111)

	// Assertions
	assert.NotNil(t, err)
	assert.False(t, deleted)
	repoProd.AssertExpectations(t)
}
