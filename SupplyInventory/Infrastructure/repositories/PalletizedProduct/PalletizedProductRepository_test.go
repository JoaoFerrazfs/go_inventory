package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"

	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	"go_inventory/SupplyInventory/tests/testutils"
)

func TestPalletizedProduct_AddProductsToPallet_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repo := &PalletizedProductRepositoryImpl{db: adapter, palletRepository: &stubPalletRepoShim{}}

	// Expectations
	// fakeAdapter and stubPalletRepoShim return success by default

	// Actions
	product := entities.PalletizedProductEntity{EAN: 999, Quantity: 1, PalletID: 1}
	ok, appErr := repo.AddProductsToPallet(product)

	// Assertions
	assert.Nil(t, appErr)
	assert.True(t, ok)
}

func TestPalletizedProduct_DeleteProductsFromPallet_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	// create repo with a stub pallet that has one product with EAN 777
	repo := &PalletizedProductRepositoryImpl{db: adapter, palletRepository: &stubPalletRepoShim{}}

	// Expectations
	// fakeAdapter.DeleteByID will return rows>0

	// Actions
	ok, appErr := repo.DeleteProductsFromPallet(1, 777)

	// Assertions
	assert.Nil(t, appErr)
	assert.True(t, ok)
}

// stub to satisfy palletRepo.PalletRepository used by the repository under test
type stubPalletRepoShim struct{}

func (s *stubPalletRepoShim) Create(p *entities.PalletEntity) error            { return nil }
func (s *stubPalletRepoShim) FindByID(id uint) (*entities.PalletEntity, error) { return nil, nil }
func (s *stubPalletRepoShim) List() ([]*entities.PalletEntity, error)          { return nil, nil }
func (s *stubPalletRepoShim) DeleteByID(id uint) error                         { return nil }
func (s *stubPalletRepoShim) Update(p *entities.PalletEntity) error            { return nil }
func (s *stubPalletRepoShim) GetAllPallets(palletRackId *uint, productId *uint) ([]entities.PalletEntity, *errors.AppError) {
	return nil, nil
}
func (s *stubPalletRepoShim) GetSupplyById(id uint) (*entities.PalletEntity, *errors.AppError) {
	return &entities.PalletEntity{ID: id, PalletizedProduct: []entities.PalletizedProductEntity{{ID: 1, EAN: 777, Quantity: 1, PalletID: id}}}, nil
}
func (s *stubPalletRepoShim) AddSupply(name string, rackId uint) (*entities.PalletEntity, *errors.AppError) {
	return nil, nil
}
func (s *stubPalletRepoShim) UpdateSupply(p *entities.PalletEntity) (*entities.PalletEntity, *errors.AppError) {
	return nil, nil
}
func (s *stubPalletRepoShim) AddProductsToPallet(prod entities.PalletizedProductEntity) (*entities.PalletEntity, *errors.AppError) {
	return nil, nil
}
func (s *stubPalletRepoShim) DeletePalletById(id uint) (bool, *errors.AppError) { return false, nil }
