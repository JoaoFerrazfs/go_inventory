package services

import (
	"log"
	"strings"

	domain "go_inventory/SupplyInventory/Domain"
	infrastructure "go_inventory/SupplyInventory/Infrastructure"
)

type PalletService interface {
	ListPallets() ([]domain.PalletEntity, error)
	FindPalletById(id uint) (*domain.PalletEntity, error)
	CreatePallet(PalletName string, PalletRackId uint) (*domain.PalletEntity, error)
}

type palletService struct {
	repo      infrastructure.PalletRepository
	qrService QRCodeService
}

func NewPalletService(repo infrastructure.PalletRepository, qrService QRCodeService) PalletService {
	return &palletService{repo: repo, qrService: qrService}
}

func (s *palletService) ListPallets() ([]domain.PalletEntity, error) {
	pallets, err := s.repo.GetAllPallets()
	if err != nil {
		return nil, err
	}
	return pallets, nil
}

func (s *palletService) FindPalletById(id uint) (*domain.PalletEntity, error) {
	return s.repo.GetSupplyById(id)
}

func (s *palletService) CreatePallet(PalletName string, PalletRackId uint) (*domain.PalletEntity, error) {
	newPallet, err := s.repo.AddSupply(PalletName, PalletRackId)
	if err != nil {
		return nil, err
	}

	qrcodeLink, err := s.qrService.CreateQRCode(newPallet.ID)
	if err != nil {
		return nil, err
	}

	newPallet.QrCode = qrcodeLink
	newPallet.QrCodeUrl = "http://localhost:3000/" + strings.TrimPrefix(qrcodeLink, "storage/")

	palletWithQrCode, err := s.repo.UpdateSupply(newPallet)
	if err != nil {
		return nil, err
	}

	log.Printf("New pallet created with QR code: %+v", palletWithQrCode)
	return palletWithQrCode, nil
}
