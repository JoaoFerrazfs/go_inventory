package services

import (
	"io"
	"strings"
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
func TestExportPalletsToCsv_WithProducts(t *testing.T) {
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
			PalletizedProduct: []entities.PalletizedProductEntity{
				{EAN: 123456789, Quantity: 5},
				{EAN: 987654321, Quantity: 3},
			},
		},
		{
			ID:             2,
			Name:           "Pallet 2",
			PalletRackID:   20,
			PalletRackName: "Rack B",
			QrCodeUrl:      "http://example.com/qr2",
			PalletizedProduct: []entities.PalletizedProductEntity{
				{EAN: 111111111, Quantity: 2},
			},
		},
		{
			ID:             3,
			Name:           "Pallet 3",
			PalletRackID:   30,
			PalletRackName: "Rack C",
			QrCodeUrl:      "http://example.com/qr3",
			PalletizedProduct: []entities.PalletizedProductEntity{}, // Pallet sem produtos
		},
	}

	expectedURL := "http://storage.example.com/reports/Pallets/test_file.csv"
	var uploadedContent string

	// Expectations
	storage.On("Upload", mock.AnythingOfType("string"), mock.MatchedBy(func(reader io.Reader) bool {
		content, _ := io.ReadAll(reader)
		uploadedContent = string(content)
		return true
	})).Return("uploaded-path", nil)
	storage.On("GetURL", mock.AnythingOfType("string")).Return(expectedURL, nil)

	// Actions
	url, err := svc.ExportPalletsToCsv(pallets)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, expectedURL, url)
	
	// Validate CSV content
	lines := strings.Split(strings.TrimSpace(uploadedContent), "\n")
	assert.Equal(t, 4, len(lines)) // Header + 3 data lines
	
	// Check header
	assert.Equal(t, "ID,Nome,Id do Rack,Nome do Rack,Produtos,QrCodeUrl", lines[0])
	
	// Check data lines
	assert.Equal(t, "1,Pallet 1,10,Rack A,\"123456789, 987654321\",http://example.com/qr1", lines[1])
	assert.Equal(t, "2,Pallet 2,20,Rack B,111111111,http://example.com/qr2", lines[2])
	assert.Equal(t, "3,Pallet 3,30,Rack C,,http://example.com/qr3", lines[3]) // Pallet sem produtos
	
	storage.AssertExpectations(t)
}