package infrastructure

import (
	errors "go_inventory/Helpers/Errors"
	repositories "go_inventory/SupplyInventory/Domain/contracts/repositories/Product"

	domain "go_inventory/SupplyInventory/Domain/Entities"

	dbadapter "go_inventory/SupplyInventory/Infrastructure/repositories/db"
)

type productRepository struct {
	db dbadapter.DBAdapter
}

func NewProductRepository(db dbadapter.DBAdapter) repositories.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ean string, name string) (*domain.ProductEntity, *errors.AppError) {
	product := domain.NewProductEntity(ean, name)

	if err := r.db.Save(&product); err != nil {
		return nil, errors.NewAppError(err.Error(), 500)
	}

	return product, nil
}

func (r *productRepository) FindByEAN(ean string) (*domain.ProductEntity, *errors.AppError) {
	var product domain.ProductEntity

	if err := r.db.WhereFirst(&product, "ean = ?", ean); err != nil {
		return nil, errors.NewAppError(err.Error(), 500)
	}

	return &product, nil
}

func (r *productRepository) Delete(ean string) (bool, *errors.AppError) {
	var product domain.ProductEntity

	if err := r.db.WhereFirst(&product, "ean = ?", ean); err != nil {
		return false, nil
	}

	_, err := r.db.DeleteByID(&product, product.ID)
	if err != nil {
		return false, errors.NewAppError(err.Error(), 500)
	}

	return true, nil
}
