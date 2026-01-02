package services

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	errors "go_inventory/Helpers/Errors"
	storage "go_inventory/SupplyInventory/Application/Services/Storage"
	entities "go_inventory/SupplyInventory/Domain/Entities"
)

type PalletExportService interface {
	ExportPalletsToCsv(pallets []entities.PalletEntity) (string, *errors.AppError)
}

type palletExportService struct {
	storage storage.Storage
}

func NewPalletExportService(storage storage.Storage) PalletExportService {
	return &palletExportService{storage: storage}
}

func (s *palletExportService) ExportPalletsToCsv(pallets []entities.PalletEntity) (string, *errors.AppError) {
	var csvBuilder strings.Builder
	writer := csv.NewWriter(&csvBuilder)
	
	// Write header
	writer.Write([]string{"ID", "Nome", "Id do Rack", "Nome do Rack", "QrCodeUrl"})

	// Write data
	for _, p := range pallets {
		writer.Write([]string{
			strconv.Itoa(int(p.ID)),
			p.Name,
			strconv.Itoa(int(p.PalletRackID)),
			p.PalletRackName,
			p.QrCodeUrl,
		})
	}

	writer.Flush()
	csvString := csvBuilder.String()

	// Generate unique filename
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	uniqueID := time.Now().UnixNano()
	filename := fmt.Sprintf("%s_%d.csv", timestamp, uniqueID)
	filePath := fmt.Sprintf("reports/Pallets/%s", filename)

	// Upload to storage
	_, err := s.storage.Upload(filePath, strings.NewReader(csvString))
	if err != nil {
		return "", errors.NewAppError(fmt.Sprintf("failed to save CSV: %v", err), 500)
	}

	// Get URL
	url, err := s.storage.GetURL(filePath)
	if err != nil {
		return "", errors.NewAppError(fmt.Sprintf("failed to get URL: %v", err), 500)
	}

	return url, nil
}