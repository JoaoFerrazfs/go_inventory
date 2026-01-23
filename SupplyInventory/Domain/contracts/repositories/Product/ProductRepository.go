package repositories

import (
	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
)

type ProductRepository interface {
	Create(ean string, name string) (*entities.ProductEntity, *errors.AppError)
	FindByEAN(ean string) (*entities.ProductEntity, *errors.AppError)
	Delete(ean string) (bool, *errors.AppError)
}
