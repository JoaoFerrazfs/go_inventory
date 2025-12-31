package services

import (
	"testing"

	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	testutils "go_inventory/SupplyInventory/tests/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPalletRepo struct{ mock.Mock }

func (m *mockPalletRepo) GetAllPallets() ([]entities.PalletEntity, *errors.AppError) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).([]entities.PalletEntity), nil
}

func (m *mockPalletRepo) GetSupplyById(id uint) (*entities.PalletEntity, *errors.AppError) {
	args := m.Called(id)
	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil
		}
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*entities.PalletEntity), nil
}

func (m *mockPalletRepo) AddSupply(name string, rackId uint) (*entities.PalletEntity, *errors.AppError) {
	args := m.Called(name, rackId)
	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil
		}
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*entities.PalletEntity), nil
}

func (m *mockPalletRepo) UpdateSupply(p *entities.PalletEntity) (*entities.PalletEntity, *errors.AppError) {
	args := m.Called(p)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*entities.PalletEntity), nil
}

func (m *mockPalletRepo) DeletePalletById(id uint) (bool, *errors.AppError) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return false, args.Get(1).(*errors.AppError)
	}
	return args.Bool(0), nil
}

// Additional methods to fully implement repositories.PalletRepository
func (m *mockPalletRepo) Create(pallet *entities.PalletEntity) error {
	args := m.Called(pallet)
	return args.Error(0)
}

func (m *mockPalletRepo) FindByID(id uint) (*entities.PalletEntity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.PalletEntity), args.Error(1)
}

func (m *mockPalletRepo) List() ([]*entities.PalletEntity, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.PalletEntity), args.Error(1)
}

func (m *mockPalletRepo) DeleteByID(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockPalletRepo) Update(pallet *entities.PalletEntity) error {
	args := m.Called(pallet)
	return args.Error(0)
}

func (m *mockPalletRepo) AddProductsToPallet(product entities.PalletizedProductEntity) (*entities.PalletEntity, *errors.AppError) {
	args := m.Called(product)
	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil
		}
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*entities.PalletEntity), nil
}

type mockQRCodeService struct{ mock.Mock }

func (m *mockQRCodeService) CreateQRCode(palletId uint) (string, string, error) {
	args := m.Called(palletId)
	return args.String(0), args.String(1), args.Error(2)
}

func TestCreatePallet_Success(t *testing.T) {
	// Set
	restore := testutils.SetEnvAndRestore("BASE_URL", "http://localhost:")
	defer restore()
	restore2 := testutils.SetEnvAndRestore("PORT", "3000")
	defer restore2()

	repo := &mockPalletRepo{}
	qr := &mockQRCodeService{}
	svc := NewPalletService(repo, qr)

	// Expectations
	created := &entities.PalletEntity{ID: 1, Name: "P1", PalletRackID: 2}
	repo.On("AddSupply", "P1", uint(2)).Return(created, nil)
	qr.On("CreateQRCode", uint(1)).Return("storage/pallet_1.png", "http://localhost:3000/pallets/1", nil)
	updated := &entities.PalletEntity{ID: 1, Name: "P1", PalletRackID: 2, QrCode: "storage/pallet_1.png", QrCodeUrl: "http://localhost:3000/pallets/1"}
	repo.On("UpdateSupply", mock.Anything).Return(updated, nil)

	// Actions
	res, err := svc.CreatePallet("P1", 2)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, updated, res)
	repo.AssertExpectations(t)
	qr.AssertExpectations(t)
}

func TestCreatePallet_QRServiceError(t *testing.T) {
	// Set
	restore := testutils.SetEnvAndRestore("BASE_URL", "http://localhost:")
	defer restore()
	restore2 := testutils.SetEnvAndRestore("PORT", "3000")
	defer restore2()

	repo := &mockPalletRepo{}
	qr := &mockQRCodeService{}
	svc := NewPalletService(repo, qr)

	// Expectations
	created := &entities.PalletEntity{ID: 2, Name: "P2", PalletRackID: 3}
	repo.On("AddSupply", "P2", uint(3)).Return(created, nil)
	qr.On("CreateQRCode", uint(2)).Return("", "", assert.AnError)

	// Actions
	res, err := svc.CreatePallet("P2", 3)

	// Assertions
	assert.Nil(t, res)
	assert.NotNil(t, err)
	repo.AssertExpectations(t)
	qr.AssertExpectations(t)
}

func TestListPallets_RepoError(t *testing.T) {
	// Set
	repo := &mockPalletRepo{}
	qr := &mockQRCodeService{}
	svc := NewPalletService(repo, qr)

	// Expectations
	appErr := errors.NewAppError("db error", 500)
	repo.On("GetAllPallets").Return(nil, appErr)

	// Actions
	res, err := svc.ListPallets()

	// Assertions
	assert.Nil(t, res)
	assert.Equal(t, appErr, err)
	repo.AssertExpectations(t)
}

func TestFindPalletById_NotFound(t *testing.T) {
	// Set
	repo := &mockPalletRepo{}
	qr := &mockQRCodeService{}
	svc := NewPalletService(repo, qr)

	// Expectations
	repo.On("GetSupplyById", uint(99)).Return(nil, nil)

	// Actions
	res, err := svc.FindPalletById(99)

	// Assertions
	assert.Nil(t, res)
	assert.Nil(t, err)
	repo.AssertExpectations(t)
}
