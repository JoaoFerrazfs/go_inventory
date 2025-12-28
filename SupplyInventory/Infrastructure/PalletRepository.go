package infrastructure

import (
	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	repositories "go_inventory/SupplyInventory/Domain/contracts/repositories/Pallet"

	"gorm.io/gorm"
)

type palletRepository struct {
	db *gorm.DB
}

// Implementação dos métodos obrigatórios da interface repositories.PalletRepository
func (repository *palletRepository) Create(pallet *entities.PalletEntity) error {
	return repository.db.Create(pallet).Error
}

func (repository *palletRepository) FindByID(id uint) (*entities.PalletEntity, error) {
	var pallet entities.PalletEntity
	if err := repository.db.First(&pallet, id).Error; err != nil {
		return nil, err
	}
	return &pallet, nil
}

func (repository *palletRepository) List() ([]*entities.PalletEntity, error) {
	var pallets []*entities.PalletEntity
	if err := repository.db.Find(&pallets).Error; err != nil {
		return nil, err
	}
	return pallets, nil
}

func (repository *palletRepository) DeleteByID(id uint) error {
	return repository.db.Delete(&entities.PalletEntity{}, id).Error
}

func (repository *palletRepository) Update(pallet *entities.PalletEntity) error {
	return repository.db.Save(pallet).Error
}

func NewPalletRepository(db *gorm.DB) repositories.PalletRepository {
	return &palletRepository{db: db}
}

func (repository *palletRepository) GetAllPallets() ([]entities.PalletEntity, *errors.AppError) {
	var pallets []entities.PalletEntity
	if err := repository.db.Preload("PalletizedProduct").Find(&pallets).Error; err != nil {
		return nil, errors.NewAppError("Pallets not found", 404)
	}
	return pallets, nil
}

func (repository *palletRepository) GetSupplyById(id uint) (*entities.PalletEntity, *errors.AppError) {
	var pallet entities.PalletEntity
	if err := repository.db.Preload("PalletizedProduct").First(&pallet, id).Error; err != nil {
		return nil, errors.NewAppError(err.Error(), 500)
	}

	return &pallet, nil
}

func (repository *palletRepository) AddSupply(PalletName string, PalletRackId uint) (*entities.PalletEntity, *errors.AppError) {
	pallet := entities.PalletEntity{
		Name:         PalletName,
		PalletRackID: PalletRackId,
	}

	if err := repository.db.Create(&pallet).Error; err != nil {
		return nil, errors.NewAppError(err.Error(), 500)
	}

	return &pallet, nil
}

func (repository *palletRepository) UpdateSupply(pallet *entities.PalletEntity) (*entities.PalletEntity, *errors.AppError) {
	if err := repository.db.Save(pallet).Error; err != nil {
		return nil, errors.NewAppError(err.Error(), 500)
	}
	return pallet, nil
}

func (repository *palletRepository) AddProductsToPallet(product entities.PalletizedProductEntity) (*entities.PalletEntity, *errors.AppError) {
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
	result := repository.db.Select("PalletizedProduct").Delete(&entities.PalletEntity{}, id)

	if result.Error != nil {
		return false, errors.NewAppError(result.Error.Error(), 500)
	}

	if result.RowsAffected == 0 {
		return false, errors.NewAppError("Pallet not found", 404)
	}

	return true, nil
}

func (repository *palletRepository) UpdatePallet(id uint, Name string, PalletRackId uint) (*entities.PalletEntity, *errors.AppError) {
	var pallet entities.PalletEntity
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

var pallet entities.PalletEntity
