package infrastructure

import (
	errors "go_inventory/Helpers/Errors"
	domain "go_inventory/SupplyInventory/Domain"
)

type PalletRepository interface {
	Create(pallet *domain.PalletEntity) error
	FindByID(id uint) (*domain.PalletEntity, error)
	List() ([]*domain.PalletEntity, error)
	DeleteByID(id uint) error
	Update(pallet *domain.PalletEntity) error
	AddProductsToPallet(product domain.PalletizedProductEntity) (*domain.PalletEntity, *errors.AppError)
	GetAllPallets() ([]domain.PalletEntity, *errors.AppError)
	GetSupplyById(id uint) (*domain.PalletEntity, *errors.AppError)
	AddSupply(name string, rackId uint) (*domain.PalletEntity, *errors.AppError)
	UpdateSupply(pallet *domain.PalletEntity) (*domain.PalletEntity, *errors.AppError)
	DeletePalletById(id uint) (bool, *errors.AppError)
}
