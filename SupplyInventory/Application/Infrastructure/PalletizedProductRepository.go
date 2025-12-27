package infrastructure

import (
	errors "go_inventory/Helpers/Errors"
	domain "go_inventory/SupplyInventory/Domain"
)

type PalletizedProductRepository interface {
	AddProductsToPallet(product domain.PalletizedProductEntity) (bool, *errors.AppError)
	DeleteProductsFromPallet(palletId uint, productsEan int) (bool, *errors.AppError)
}
