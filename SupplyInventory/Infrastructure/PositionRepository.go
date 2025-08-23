package infrastructure

import (
	"log"

	domain "go_inventory/SupplyInventory/Domain"

	"gorm.io/gorm"
)

// Interface para facilitar testes e desacoplamento
type PalletRepository interface {
	GetAllPallets() ([]domain.Pallet, error)
	GetSupplyById(id uint) (*domain.Pallet, error)
	AddSupply(pallet domain.Pallet) (*domain.Pallet, error)
	UpdateSupply(pallet *domain.Pallet) (*domain.Pallet, error)
	AddProductsToPallet(product domain.PalletizedProductEntity) (*domain.Pallet, error)
}

// Implementação concreta
type palletRepository struct {
	db *gorm.DB
}

// Construtor
func NewPalletRepository(db *gorm.DB) PalletRepository {
	return &palletRepository{db: db}
}

func (r *palletRepository) GetAllPallets() ([]domain.Pallet, error) {
	var pallets []domain.Pallet
	if err := r.db.Preload("PalletizedProduct").Find(&pallets).Error; err != nil {
		return nil, err
	}
	return pallets, nil
}

func (r *palletRepository) GetSupplyById(id uint) (*domain.Pallet, error) {
	var pallet domain.Pallet
	if err := r.db.Preload("PalletizedProduct").First(&pallet, id).Error; err != nil {
		return nil, err
	}
	log.Printf("Pallet found: %+v", pallet)
	return &pallet, nil
}

func (r *palletRepository) AddSupply(pallet domain.Pallet) (*domain.Pallet, error) {
	if err := r.db.Create(&pallet).Error; err != nil {
		return nil, err
	}
	return &pallet, nil
}

func (r *palletRepository) UpdateSupply(pallet *domain.Pallet) (*domain.Pallet, error) {
	if err := r.db.Save(pallet).Error; err != nil {
		return nil, err
	}
	return pallet, nil
}

func (r *palletRepository) AddProductsToPallet(product domain.PalletizedProductEntity) (*domain.Pallet, error) {
	pallet, err := r.GetSupplyById(product.PalletID)
	if err != nil || pallet == nil {
		log.Print(product.PalletID, 5)
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
