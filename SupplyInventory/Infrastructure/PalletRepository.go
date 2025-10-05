package infrastructure

import (
	errors "go_inventory/Helpers/Errors"
	domain "go_inventory/SupplyInventory/Domain"

	"gorm.io/gorm"
)

type PalletRepository interface {
	GetAllPallets() ([]domain.PalletEntity, *errors.AppError)
	GetSupplyById(id uint) (*domain.PalletEntity, *errors.AppError)
	AddSupply(PalletName string, PalletRackId uint) (*domain.PalletEntity, *errors.AppError)
	UpdateSupply(pallet *domain.PalletEntity) (*domain.PalletEntity, *errors.AppError)
	AddProductsToPallet(product domain.PalletizedProductEntity) (*domain.PalletEntity, *errors.AppError)
	DeletePalletById(id uint) (bool, *errors.AppError)
	UpdatePallet(id uint, Name string, PalletRackId uint) (*domain.PalletEntity, *errors.AppError)
}

type palletRepository struct {
	db *gorm.DB
}

func NewPalletRepository(db *gorm.DB) PalletRepository {
	return &palletRepository{db: db}
}

func (repository *palletRepository) GetAllPallets() ([]domain.PalletEntity, *errors.AppError) {
	var pallets []domain.PalletEntity
	if err := repository.db.Preload("PalletizedProduct").Find(&pallets).Error; err != nil {
		return nil, errors.NewAppError("Pallets not found", 404)
	}
	return pallets, nil
}

func (repository *palletRepository) GetSupplyById(id uint) (*domain.PalletEntity, *errors.AppError) {
	var pallet domain.PalletEntity
	if err := repository.db.Preload("PalletizedProduct").First(&pallet, id).Error; err != nil {
		return nil, errors.NewAppError(err.Error(), 500)
	}

	return &pallet, nil
}

func (repository *palletRepository) AddSupply(PalletName string, PalletRackId uint) (*domain.PalletEntity, *errors.AppError) {
	pallet := domain.PalletEntity{
		Name:         PalletName,
		PalletRackID: PalletRackId,
	}

	if err := repository.db.Create(&pallet).Error; err != nil {
		return nil, errors.NewAppError(err.Error(), 500)
	}

	return &pallet, nil
}

func (repository *palletRepository) UpdateSupply(pallet *domain.PalletEntity) (*domain.PalletEntity, *errors.AppError) {
	if err := repository.db.Save(pallet).Error; err != nil {
		return nil, errors.NewAppError(err.Error(), 500)
	}
	return pallet, nil
}

func (repository *palletRepository) AddProductsToPallet(product domain.PalletizedProductEntity) (*domain.PalletEntity, *errors.AppError) {
	pallet, err := repository.GetSupplyById(product.PalletID)
	if err != nil || pallet == nil {
		return nil, errors.NewAppError("Pallet not found", 404)
	}

	product.PalletID = pallet.ID
	if err := repository.db.Model(pallet).Association("PalletizedProduct").Append(&product); err != nil {
		return nil, errors.NewAppError(err.Error(), 422)
	}

	if err := repository.db.Preload("PalletizedProduct").First(pallet, pallet.ID).Error; err != nil {
		return nil, errors.NewAppError(err.Error(), 400)
	}

	return pallet, nil
}

func (repository *palletRepository) DeletePalletById(id uint) (bool, *errors.AppError) {
	result := repository.db.Select("PalletizedProduct").Delete(&domain.PalletEntity{}, id)

	if result.Error != nil {
		return false, errors.NewAppError(result.Error.Error(), 500)
	}

	if result.RowsAffected == 0 {
		return false, errors.NewAppError("Pallet not found", 404)
	}

	return true, nil
}

func (repository *palletRepository) UpdatePallet(id uint, Name string, PalletRackId uint) (*domain.PalletEntity, *errors.AppError) {
	var pallet domain.PalletEntity
	if err := repository.db.Preload("PalletizedProduct").First(&pallet, id).Error; err != nil {
		return nil, errors.NewAppError("Pallet not found", 404)
	}

	pallet.Name = Name
	pallet.PalletRackID = PalletRackId

	if err := repository.db.Save(&pallet).Error; err != nil {
		return nil, errors.NewAppError(err.Error(), 500)
	}

	return &pallet, nil
}
