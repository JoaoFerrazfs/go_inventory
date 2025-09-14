package services

import (
	"strings"

	errors "go_inventory/Helpers/Errors"
	domain "go_inventory/SupplyInventory/Domain"
	infrastructure "go_inventory/SupplyInventory/Infrastructure"
)

type PalletService interface {
	ListPallets() ([]domain.PalletEntity, *errors.AppError)
	FindPalletById(id uint) (*domain.PalletEntity, *errors.AppError)
	CreatePallet(PalletName string, PalletRackId uint) (*domain.PalletEntity, *errors.AppError)
	DeletePalletById(id uint) (bool, *errors.AppError)
}

type palletService struct {
	repo      infrastructure.PalletRepository
	qrService QRCodeService
}

func NewPalletService(repo infrastructure.PalletRepository, qrService QRCodeService) PalletService {
	return &palletService{repo: repo, qrService: qrService}
}

func (s *palletService) ListPallets() ([]domain.PalletEntity, *errors.AppError) {
	pallets, appErr := s.repo.GetAllPallets()
	if appErr != nil {
		return nil, appErr
	}
	return pallets, nil
}

func (s *palletService) FindPalletById(id uint) (*domain.PalletEntity, *errors.AppError) {
	return s.repo.GetSupplyById(id)
}

func (s *palletService) CreatePallet(PalletName string, PalletRackId uint) (*domain.PalletEntity, *errors.AppError) {
	newPallet, appErr := s.repo.AddSupply(PalletName, PalletRackId)
	if appErr != nil {
		return nil, appErr
	}

	qrcodeLink, err := s.qrService.CreateQRCode(newPallet.ID)
	if err != nil {
		return nil, errors.NewAppError(err.Error(), 404)
	}

	newPallet.QrCode = qrcodeLink
	newPallet.QrCodeUrl = "http://localhost:3000/" + strings.TrimPrefix(qrcodeLink, "storage/")

	palletWithQrCode, appErr := s.repo.UpdateSupply(newPallet)
	if appErr != nil {
		return nil, appErr
	}

	return palletWithQrCode, nil
}

func (s *palletService) DeletePalletById(id uint) (bool, *errors.AppError) {
	return s.repo.DeletePalletById(id)
}
