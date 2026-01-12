package services_test

import (
	"testing"

	errors "go_inventory/Helpers/Errors"
	palletService "go_inventory/SupplyInventory/Application/Services/Pallet"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	"go_inventory/SupplyInventory/tests/mocks"
	testutils "go_inventory/SupplyInventory/tests/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPalletRepo struct{ mock.Mock }

func (m *mockPalletRepo) GetAllPallets(palletRackId *uint, productEan *int) ([]entities.PalletEntity, *errors.AppError) {
	args := m.Called(palletRackId, productEan)
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

type mockPalletRackRepo struct{ mock.Mock }

func (m *mockPalletRackRepo) Create(name string, location string, totalCapacity int, inventoryID uint) (*entities.PalletRackEntity, error) {
	args := m.Called(name, location, totalCapacity, inventoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.PalletRackEntity), args.Error(1)
}

func (m *mockPalletRackRepo) ListRacks(inventoryID uint) ([]entities.PalletRackEntity, error) {
	args := m.Called(inventoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.PalletRackEntity), args.Error(1)
}

func (m *mockPalletRackRepo) FindPalletById(id uint) (*entities.PalletRackEntity, *errors.AppError) {
	args := m.Called(id)
	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil
		}
		return nil, args.Get(1).(*errors.AppError)
	}
	return args.Get(0).(*entities.PalletRackEntity), nil
}

func (m *mockPalletRackRepo) DeleteRack(id uint) (bool, *errors.AppError) {
	args := m.Called(id)
	if args.Get(1) != nil {
		return args.Bool(0), args.Get(1).(*errors.AppError)
	}
	return args.Bool(0), nil
}

type mockPalletExportService struct{ mock.Mock }

func (m *mockPalletExportService) ExportPalletsToCsv(pallets []entities.PalletEntity) (string, *errors.AppError) {
	args := m.Called(pallets)
	return args.String(0), args.Get(1).(*errors.AppError)
}

func TestCreatePallet_Success(t *testing.T) {
	// Set
	restore := testutils.SetEnvAndRestore("BASE_URL", "http://localhost:")
	defer restore()
	restore2 := testutils.SetEnvAndRestore("PORT", "3000")
	defer restore2()

	palletRepository := &mockPalletRepo{}
	qrCodeService := &mockQRCodeService{}
	palletRackRepository := &mockPalletRackRepo{}
	storageService := &mocks.MockStorage{}
	palletExportService := &mockPalletExportService{}
	palletService := palletService.NewPalletService(palletRepository, qrCodeService, palletRackRepository, storageService, palletExportService)

	// Expectations
	rack := &entities.PalletRackEntity{ID: 2, Name: "Rack1", InventoryID: testutils.DefaultInventoryID}
	palletRackRepository.On("FindPalletById", uint(2)).Return(rack, nil)
	newPallet := &entities.PalletEntity{ID: 1, Name: "P1", PalletRackID: 2, InventoryID: testutils.DefaultInventoryID}
	palletRepository.On("AddSupply", "P1", uint(2)).Return(newPallet, nil)
	qrCodeService.On("CreateQRCode", uint(1)).Return("storage/pallet_1.png", "http://localhost:3000/pallets/1", nil)
	updated := &entities.PalletEntity{ID: 1, Name: "P1", PalletRackID: 2, PalletRackName: "Rack1", QrCode: "storage/pallet_1.png", QrCodeUrl: "http://localhost:3000/pallets/1", InventoryID: testutils.DefaultInventoryID}
	palletRepository.On("UpdateSupply", mock.Anything).Return(updated, nil)

	// Actions
	res, err := palletService.CreatePallet("P1", 2, testutils.DefaultInventoryID)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, updated, res)
	palletRepository.AssertExpectations(t)
	qrCodeService.AssertExpectations(t)
	palletRackRepository.AssertExpectations(t)
}

func TestCreatePallet_QRServiceError(t *testing.T) {
	// Set
	restore := testutils.SetEnvAndRestore("BASE_URL", "http://localhost:")
	defer restore()
	restorePort := testutils.SetEnvAndRestore("PORT", "3000")
	defer restorePort()

	palletRepository := &mockPalletRepo{}
	qrCodeService := &mockQRCodeService{}
	palletRackRepository := &mockPalletRackRepo{}
	storageService := &mocks.MockStorage{}
	palletExportService := &mockPalletExportService{}
	palletService := palletService.NewPalletService(palletRepository, qrCodeService, palletRackRepository, storageService, palletExportService)

	// Expectations
	palletRack := &entities.PalletRackEntity{ID: 3, Name: "Rack3", InventoryID: testutils.DefaultInventoryID}
	palletRackRepository.On("FindPalletById", uint(3)).Return(palletRack, nil)
	newPallet := &entities.PalletEntity{ID: 2, Name: "P2", PalletRackID: 3, InventoryID: testutils.DefaultInventoryID}
	palletRepository.On("AddSupply", "P2", uint(3)).Return(newPallet, nil)
	qrCodeService.On("CreateQRCode", uint(2)).Return("", "", assert.AnError)

	// Actions
	res, err := palletService.CreatePallet("P2", 3, testutils.DefaultInventoryID)

	// Assertions
	assert.Nil(t, res)
	assert.NotNil(t, err)
	palletRepository.AssertExpectations(t)
	qrCodeService.AssertExpectations(t)
	palletRackRepository.AssertExpectations(t)
}

func TestListPallets_RepoError(t *testing.T) {
	// Set
	palletRepository := &mockPalletRepo{}
	qrCodeService := &mockQRCodeService{}
	palletRackRepository := &mockPalletRackRepo{}
	storageService := &mocks.MockStorage{}
	palletExportService := &mockPalletExportService{}
	palletService := palletService.NewPalletService(palletRepository, qrCodeService, palletRackRepository, storageService, palletExportService)

	// Expectations
	appErr := errors.NewAppError("db error", 500)
	palletRepository.On("GetAllPallets", (*uint)(nil), (*int)(nil)).Return(nil, appErr)

	// Actions
	res, err := palletService.ListPallets(nil, nil)

	// Assertions
	assert.Nil(t, res)
	assert.Equal(t, appErr, err)
	palletRepository.AssertExpectations(t)
}

func TestFindPalletById_NotFound(t *testing.T) {
	// Set
	palletRepository := &mockPalletRepo{}
	qrCodeService := &mockQRCodeService{}
	palletRackRepository := &mockPalletRackRepo{}
	storageService := &mocks.MockStorage{}
	palletExportService := &mockPalletExportService{}
	palletService := palletService.NewPalletService(palletRepository, qrCodeService, palletRackRepository, storageService, palletExportService)

	// Expectations
	palletRepository.On("GetSupplyById", uint(99)).Return(nil, nil)

	// Actions
	res, err := palletService.FindPalletById(99)

	// Assertions
	assert.Nil(t, res)
	assert.Nil(t, err)
	palletRepository.AssertExpectations(t)
}
