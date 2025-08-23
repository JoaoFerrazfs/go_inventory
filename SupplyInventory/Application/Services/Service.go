package services

import (
	"log"
	"strings"

	domain "go_inventory/SupplyInventory/Domain"
	infrastructure "go_inventory/SupplyInventory/Infrastructure"
)

type PalletService interface {
	ListPallets() []domain.PalletEntity
	FindPalletById(id uint) (*domain.PalletEntity, error)
	CreatePallet(pallet domain.PalletEntity) (*domain.PalletEntity, error)
	AddProductsToPallet(product domain.PalletizedProductEntity) *domain.PalletEntity
}

type palletService struct {
	repo      infrastructure.PalletRepository
	qrService QRCodeService
}

func NewPalletService(repo infrastructure.PalletRepository, qrService QRCodeService) PalletService {
	return &palletService{repo: repo, qrService: qrService}
}

func (s *palletService) ListPallets() []domain.PalletEntity {
	pallets, err := s.repo.GetAllPallets()
	if err != nil {
		return nil
	}
	return pallets
}

func (s *palletService) FindPalletById(id uint) (*domain.PalletEntity, error) {
	return s.repo.GetSupplyById(id)
}

func (s *palletService) CreatePallet(pallet domain.PalletEntity) (*domain.PalletEntity, error) {
	newPallet, err := s.repo.AddSupply(pallet)
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

func (s *palletService) AddProductsToPallet(product domain.PalletizedProductEntity) *domain.PalletEntity {
	newPalletProduct, err := s.repo.AddProductsToPallet(product)
	if err != nil {
		return nil
	}
	return newPalletProduct
}
