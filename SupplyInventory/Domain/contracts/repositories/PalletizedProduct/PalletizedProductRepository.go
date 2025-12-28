package repositories

import (
	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
)

type PalletizedProductRepository interface {
    AddProductsToPallet(product entities.PalletizedProductEntity) (bool, *errors.AppError)
    DeleteProductsFromPallet(palletId uint, productsEan int) (bool, *errors.AppError)
}
