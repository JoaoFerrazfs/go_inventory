package infrastructure

import (
	domain "go_inventory/SupplyInventory/Domain"

	"gorm.io/gorm"
)

type PalletRepository interface {
	GetAllPallets() ([]domain.PalletEntity, error)
	GetSupplyById(id uint) (*domain.PalletEntity, error)
	AddSupply(pallet domain.PalletEntity) (*domain.PalletEntity, error)
	UpdateSupply(pallet *domain.PalletEntity) (*domain.PalletEntity, error)
	AddProductsToPallet(product domain.PalletizedProductEntity) (*domain.PalletEntity, error)
}

type palletRepository struct {
	db *gorm.DB
}

func NewPalletRepository(db *gorm.DB) PalletRepository {
	return &palletRepository{db: db}
}

func (r *palletRepository) GetAllPallets() ([]domain.PalletEntity, error) {
	var pallets []domain.PalletEntity
	if err := r.db.Preload("PalletizedProduct").Find(&pallets).Error; err != nil {
		return nil, err
	}
	return pallets, nil
}

func (r *palletRepository) GetSupplyById(id uint) (*domain.PalletEntity, error) {
	var pallet domain.PalletEntity
	if err := r.db.Preload("PalletizedProduct").First(&pallet, id).Error; err != nil {
		return nil, err
	}

	return &pallet, nil
}

func (r *palletRepository) AddSupply(pallet domain.PalletEntity) (*domain.PalletEntity, error) {
	if err := r.db.Create(&pallet).Error; err != nil {
		return nil, err
	}
	return &pallet, nil
}

func (r *palletRepository) UpdateSupply(pallet *domain.PalletEntity) (*domain.PalletEntity, error) {
	if err := r.db.Save(pallet).Error; err != nil {
		return nil, err
	}
	return pallet, nil
}

func (r *palletRepository) AddProductsToPallet(product domain.PalletizedProductEntity) (*domain.PalletEntity, error) {
	pallet, err := r.GetSupplyById(product.PalletID)
	if err != nil || pallet == nil {
		return nil, err
	}

	product.PalletID = pallet.ID
	if err := r.db.Model(pallet).Association("PalletizedProduct").Append(&product); err != nil {
		return nil, err
	}

	if err := r.db.Preload("PalletizedProduct").First(pallet, pallet.ID).Error; err != nil {
		return nil, err
	}

	return pallet, nil
}
