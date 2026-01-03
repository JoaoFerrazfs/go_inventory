package services

import (
	"testing"

	entities "go_inventory/SupplyInventory/Domain/Entities"
	"go_inventory/SupplyInventory/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExportPalletsToCsv_Success(t *testing.T) {
	// Set
	storage := &mocks.MockStorage{}
	svc := NewPalletExportService(storage)

	pallets := []entities.PalletEntity{
		{
			ID:             1,
			Name:           "Pallet 1",
			PalletRackID:   10,
			PalletRackName: "Rack A",
			QrCodeUrl:      "http://example.com/qr1",
		},
		{
			ID:             2,
			Name:           "Pallet 2",
			PalletRackID:   20,
			PalletRackName: "Rack B",
			QrCodeUrl:      "http://example.com/qr2",
		},
	}

	expectedURL := "http://storage.example.com/reports/Pallets/test_file.csv"

	// Expectations
	storage.On("Upload", mock.AnythingOfType("string"), mock.AnythingOfType("*strings.Reader")).Return("uploaded-path", nil)
	storage.On("GetURL", mock.AnythingOfType("string")).Return(expectedURL, nil)

	// Actions
	url, err := svc.ExportPalletsToCsv(pallets)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, expectedURL, url)
	storage.AssertExpectations(t)
}

func TestExportPalletsToCsv_UploadError(t *testing.T) {
	// Set
	storage := &mocks.MockStorage{}
	svc := NewPalletExportService(storage)

	pallets := []entities.PalletEntity{
		{
			ID:             1,
			Name:           "Pallet 1",
			PalletRackID:   10,
			PalletRackName: "Rack A",
			QrCodeUrl:      "http://example.com/qr1",
		},
	}

	// Expectations
	storage.On("Upload", mock.AnythingOfType("string"), mock.AnythingOfType("*strings.Reader")).Return("", assert.AnError)

	// Actions
	url, err := svc.ExportPalletsToCsv(pallets)

	// Assertions
	assert.Error(t, err)
	assert.Empty(t, url)
	assert.Contains(t, err.Message, "failed to save CSV")
	assert.Equal(t, 500, err.Code)
	storage.AssertExpectations(t)
}

func TestExportPalletsToCsv_GetURLError(t *testing.T) {
	// Set
	storage := &mocks.MockStorage{}
	svc := NewPalletExportService(storage)

	pallets := []entities.PalletEntity{
		{
			ID:             1,
			Name:           "Pallet 1",
			PalletRackID:   10,
			PalletRackName: "Rack A",
			QrCodeUrl:      "http://example.com/qr1",
		},
	}

	// Expectations
	storage.On("Upload", mock.AnythingOfType("string"), mock.AnythingOfType("*strings.Reader")).Return("uploaded-path", nil)
	storage.On("GetURL", mock.AnythingOfType("string")).Return("", assert.AnError)

	// Actions
	url, err := svc.ExportPalletsToCsv(pallets)

	// Assertions
	assert.Error(t, err)
	assert.Empty(t, url)
	assert.Contains(t, err.Message, "failed to get URL")
	assert.Equal(t, 500, err.Code)
	storage.AssertExpectations(t)
}

func TestExportPalletsToCsv_EmptyPallets(t *testing.T) {
	// Set
	storage := &mocks.MockStorage{}
	svc := NewPalletExportService(storage)

	pallets := []entities.PalletEntity{}
	expectedURL := "http://storage.example.com/reports/Pallets/empty_file.csv"

	// Expectations
	storage.On("Upload", mock.AnythingOfType("string"), mock.Anything).Return("uploaded-path", nil)
	storage.On("GetURL", mock.AnythingOfType("string")).Return(expectedURL, nil)

	// Actions
	url, err := svc.ExportPalletsToCsv(pallets)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, expectedURL, url)
	storage.AssertExpectations(t)
}
