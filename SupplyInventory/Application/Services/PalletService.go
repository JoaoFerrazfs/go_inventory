package services

import (
	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	repositories "go_inventory/SupplyInventory/Domain/contracts/repositories/Pallet"
	"strings"
)

type PalletService interface {
	ListPallets() ([]entities.PalletEntity, *errors.AppError)
	FindPalletById(id uint) (*entities.PalletEntity, *errors.AppError)
	CreatePallet(PalletName string, PalletRackId uint) (*entities.PalletEntity, *errors.AppError)
	DeletePalletById(id uint) (bool, *errors.AppError)
	UpdatePallet(id uint, Name string, PalletRackId uint) (*entities.PalletEntity, *errors.AppError)
}

type palletService struct {
	repo      repositories.PalletRepository
	qrService QRCodeService
}

func NewPalletService(repo repositories.PalletRepository, qrService QRCodeService) PalletService {
	return &palletService{repo: repo, qrService: qrService}
}

func (service *palletService) ListPallets() ([]entities.PalletEntity, *errors.AppError) {
	pallets, appErr := service.repo.GetAllPallets()
	if appErr != nil {
		return nil, appErr
	}
	return pallets, nil
}

func (service *palletService) FindPalletById(id uint) (*entities.PalletEntity, *errors.AppError) {
	return service.repo.GetSupplyById(id)
}

func (service *palletService) CreatePallet(PalletName string, PalletRackId uint) (*entities.PalletEntity, *errors.AppError) {
	newPallet, appErr := service.repo.AddSupply(PalletName, PalletRackId)
	if appErr != nil {
		return nil, appErr
	}

	qrcodeLink, err := service.qrService.CreateQRCode(newPallet.ID)
	if err != nil {
		return nil, errors.NewAppError(err.Error(), 404)
	}

	newPallet.QrCode = qrcodeLink
	newPallet.QrCodeUrl = "http://localhost:3000/" + strings.TrimPrefix(qrcodeLink, "storage/")

	palletWithQrCode, appErr := service.repo.UpdateSupply(newPallet)
	if appErr != nil {
		return nil, appErr
	}

	return palletWithQrCode, nil
}

func (service *palletService) DeletePalletById(id uint) (bool, *errors.AppError) {
	return service.repo.DeletePalletById(id)
}

func (service *palletService) UpdatePallet(id uint, Name string, PalletRackId uint) (*entities.PalletEntity, *errors.AppError) {
	pallet := &entities.PalletEntity{
		ID: id,
		Name: Name,
		PalletRackID: PalletRackId,
	}
	return service.repo.UpdateSupply(pallet)
}
