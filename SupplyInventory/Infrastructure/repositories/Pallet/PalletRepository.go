package infrastructure

import (
	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	repositories "go_inventory/SupplyInventory/Domain/contracts/repositories/Pallet"
	dbadapter "go_inventory/SupplyInventory/Infrastructure/repositories/db"
)

type palletRepository struct {
	db dbadapter.DBAdapter
}

// Implementação dos métodos obrigatórios da interface repositories.PalletRepository
func (repository *palletRepository) Create(pallet *entities.PalletEntity) error {
	return repository.db.Create(pallet)
}

func (repository *palletRepository) FindByID(id uint) (*entities.PalletEntity, error) {
	var pallet entities.PalletEntity
	if err := repository.db.FirstByID(&pallet, id); err != nil {
		return nil, err
	}
	return &pallet, nil
}

func (repository *palletRepository) List() ([]*entities.PalletEntity, error) {
	var pallets []*entities.PalletEntity
	if err := repository.db.FindAll(&pallets); err != nil {
		return nil, err
	}
	return pallets, nil
}

func (repository *palletRepository) DeleteByID(id uint) error {
	_, err := repository.db.DeleteByID(&entities.PalletEntity{}, id)
	return err
}

func (repository *palletRepository) Update(pallet *entities.PalletEntity) error {
	return repository.db.Save(pallet)
}

func NewPalletRepository(db dbadapter.DBAdapter) repositories.PalletRepository {
	return &palletRepository{db: db}
}

func (repository *palletRepository) GetAllPallets(palletRackId *uint, productId *uint) ([]entities.PalletEntity, *errors.AppError) {
	var pallets []entities.PalletEntity
	query := repository.db.GetDB().Preload("PalletizedProduct")
	
	if palletRackId != nil {
		query = query.Where("pallet_rack_id = ?", *palletRackId)
	}
	
	if productId != nil {
		query = query.Joins("JOIN palletized_products pp ON pp.pallet_id = pallets.id").Where("pp.id = ?", *productId)
	}
	
	if err := query.Find(&pallets).Error; err != nil {
		return nil, errors.NewAppError("Pallets not found", 404)
	}
	return pallets, nil
}

func (repository *palletRepository) GetSupplyById(id uint) (*entities.PalletEntity, *errors.AppError) {
	var pallet entities.PalletEntity
	if err := repository.db.PreloadFind(&pallet, "PalletizedProduct", id); err != nil {
		return nil, errors.NewAppError(err.Error(), 500)
	}

	return &pallet, nil
}

func (repository *palletRepository) AddSupply(PalletName string, PalletRackId uint) (*entities.PalletEntity, *errors.AppError) {
	pallet := entities.PalletEntity{
		Name:         PalletName,
		PalletRackID: PalletRackId,
	}

	if err := repository.db.Create(&pallet); err != nil {
		return nil, errors.NewAppError(err.Error(), 500)
	}

	return &pallet, nil
}

func (repository *palletRepository) UpdateSupply(pallet *entities.PalletEntity) (*entities.PalletEntity, *errors.AppError) {
	if err := repository.db.Save(pallet); err != nil {
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
	if err := repository.db.AppendAssociation(pallet, &product); err != nil {
		return nil, errors.NewAppError(err.Error(), 422)
	}

	if err := repository.db.PreloadFind(pallet, "PalletizedProduct", pallet.ID); err != nil {
		return nil, errors.NewAppError(err.Error(), 400)
	}

	return pallet, nil
}

func (repository *palletRepository) DeletePalletById(id uint) (bool, *errors.AppError) {
	rows, err := repository.db.DeleteByID(&entities.PalletEntity{}, id)
	if err != nil {
		return false, errors.NewAppError(err.Error(), 500)
	}
	if rows == 0 {
		return false, errors.NewAppError("Pallet not found", 404)
	}
	return true, nil
}
