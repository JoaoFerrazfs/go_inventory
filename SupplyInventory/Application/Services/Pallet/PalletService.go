package services

import (
	errors "go_inventory/Helpers/Errors"
	qrCodeService "go_inventory/SupplyInventory/Application/Services/QrCode"
	storage "go_inventory/SupplyInventory/Application/Services/Storage"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	palletRepository "go_inventory/SupplyInventory/Domain/contracts/repositories/Pallet"
	palletRackRepository "go_inventory/SupplyInventory/Domain/contracts/repositories/PalletRack"
)

type PalletService interface {
	ListPallets(palletRackId *uint, productEan *int) ([]entities.PalletEntity, *errors.AppError)
	FindPalletById(id uint) (*entities.PalletEntity, *errors.AppError)
	CreatePallet(PalletName string, PalletRackId uint) (*entities.PalletEntity, *errors.AppError)
	DeletePalletById(id uint) (bool, *errors.AppError)
	UpdatePallet(id uint, Name string, PalletRackId uint) (*entities.PalletEntity, *errors.AppError)
	GeneratePalletsCsvFile(palletRackId *uint, productEan *int) (string, *errors.AppError)
}

type palletService struct {
	palletRepository     palletRepository.PalletRepository
	qrService            qrCodeService.QRCodeService
	palletRackRepository palletRackRepository.PalletRackRepository
	storage              storage.Storage
	exportService        PalletExportService
}

func NewPalletService(palletRepository palletRepository.PalletRepository, qrService qrCodeService.QRCodeService, palletRackRepository palletRackRepository.PalletRackRepository, storage storage.Storage, exportService PalletExportService) PalletService {
	return &palletService{palletRepository: palletRepository, qrService: qrService, palletRackRepository: palletRackRepository, storage: storage, exportService: exportService}
}

func (service *palletService) ListPallets(palletRackId *uint, productEan *int) ([]entities.PalletEntity, *errors.AppError) {
	pallets, appErr := service.palletRepository.GetAllPallets(palletRackId, productEan)
	if appErr != nil {
		return nil, appErr
	}
	return pallets, nil
}

func (service *palletService) FindPalletById(id uint) (*entities.PalletEntity, *errors.AppError) {
	return service.palletRepository.GetSupplyById(id)
}

func (service *palletService) CreatePallet(PalletName string, PalletRackId uint) (*entities.PalletEntity, *errors.AppError) {

	palletRack, apperr := service.palletRackRepository.FindPalletById(PalletRackId)
	if apperr != nil {
		return nil, apperr
	}

	newPallet, appErr := service.palletRepository.AddSupply(PalletName, PalletRackId)
	if appErr != nil {
		return nil, appErr
	}

	storagePath, publicURL, err := service.qrService.CreateQRCode(newPallet.ID)
	if err != nil {
		return nil, errors.NewAppError(err.Error(), 404)
	}

	newPallet.QrCode = storagePath
	newPallet.QrCodeUrl = publicURL
	newPallet.PalletRackName = palletRack.Name
	newPallet.PalletRackID = palletRack.ID

	palletWithQrCode, appErr := service.palletRepository.UpdateSupply(newPallet)
	if appErr != nil {
		return nil, appErr
	}

	return palletWithQrCode, nil
}

func (service *palletService) DeletePalletById(id uint) (bool, *errors.AppError) {
	return service.palletRepository.DeletePalletById(id)
}

func (service *palletService) UpdatePallet(id uint, Name string, PalletRackId uint) (*entities.PalletEntity, *errors.AppError) {
	// Load existing pallet to avoid overwriting other fields
	existing, appErr := service.palletRepository.GetSupplyById(id)
	if appErr != nil {
		return nil, appErr
	}

	palletRack, apperr := service.palletRackRepository.FindPalletById(PalletRackId)
	if apperr != nil {
		return nil, apperr
	}

	existing.Name = Name
	existing.PalletRackID = PalletRackId
	existing.PalletRackName = palletRack.Name

	return service.palletRepository.UpdateSupply(existing)
}

func (service *palletService) GeneratePalletsCsvFile(palletRackId *uint, productEan *int) (string, *errors.AppError) {
	pallets, appErr := service.ListPallets(palletRackId, productEan)
	if appErr != nil {
		return "", appErr
	}

	return service.exportService.ExportPalletsToCsv(pallets)
}
