package repositories

import (
	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
)

type PalletRepository interface {
    Create(pallet *entities.PalletEntity) error
    FindByID(id uint) (*entities.PalletEntity, error)
    List() ([]*entities.PalletEntity, error)
    DeleteByID(id uint) error
    Update(pallet *entities.PalletEntity) error
    AddProductsToPallet(product entities.PalletizedProductEntity) (*entities.PalletEntity, *errors.AppError)
    GetAllPallets() ([]entities.PalletEntity, *errors.AppError)
    GetSupplyById(id uint) (*entities.PalletEntity, *errors.AppError)
    AddSupply(name string, rackId uint) (*entities.PalletEntity, *errors.AppError)
    UpdateSupply(pallet *entities.PalletEntity) (*entities.PalletEntity, *errors.AppError)
    DeletePalletById(id uint) (bool, *errors.AppError)
}
