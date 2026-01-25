package services

import (
	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	repositories "go_inventory/SupplyInventory/Domain/contracts/repositories/Product"
)

type ProductService interface {
	CreateProduct(ean string, name string) (*entities.ProductEntity, *errors.AppError)
	GetProductByEAN(ean string) (*entities.ProductEntity, *errors.AppError)
	DeleteProduct(ean string) (bool, *errors.AppError)
}

type productService struct {
	productRepository repositories.ProductRepository
}

func NewProductService(productRepository repositories.ProductRepository) ProductService {
	return &productService{
		productRepository: productRepository,
	}
}

func (service *productService) CreateProduct(ean string, name string) (*entities.ProductEntity, *errors.AppError) {
	existingProduct, _ := service.productRepository.FindByEAN(ean)
	if existingProduct != nil {
		return nil, errors.NewAppError("Product with given EAN already exists", 400)
	}

	product, appErr := service.productRepository.Create(ean, name)
	if appErr != nil {
		return nil, appErr
	}

	return product, nil
}

func (service *productService) GetProductByEAN(ean string) (*entities.ProductEntity, *errors.AppError) {
	product, appErr := service.productRepository.FindByEAN(ean)
	if appErr != nil {
		return nil, appErr
	}

	return product, nil
}

func (service *productService) DeleteProduct(ean string) (bool, *errors.AppError) {
	deleted, appErr := service.productRepository.Delete(ean)
	if appErr != nil {
		return false, appErr
	}

	return deleted, nil
}
