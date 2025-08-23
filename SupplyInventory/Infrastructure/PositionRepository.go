package infrastructure

import (
	"log"

	domain "go_inventory/SupplyInventory/Domain"
	db "go_inventory/SupplyInventory/Infrastructure/Db"
)

func GetAllPallets() ([]domain.Pallet, error) {
	var pallets []domain.Pallet
	if err := db.DB.Preload("PalletizedProduct").Find(&pallets).Error; err != nil {
		return nil, err
	}
	return pallets, nil
}

func GetSupplyById(id uint) (*domain.Pallet, error) {
	var pallet domain.Pallet
	if err := db.DB.Preload("PalletizedProduct").First(&pallet, id).Error; err != nil {
		return nil, err
	}
	log.Printf("Pallet found: %+v", pallet)
	return &pallet, nil
}

func AddSupply(pallet domain.Pallet) (*domain.Pallet, error) {
	if err := db.DB.Create(&pallet).Error; err != nil {
		return nil, err
	}
	return &pallet, nil
}

func UpdateSupply(pallet *domain.Pallet) (*domain.Pallet, error) {
	if err := db.DB.Save(&pallet).Error; err != nil {
		return nil, err
	}
	return pallet, nil
}

func AddProductsToPallet(product domain.PalletizedProductEntity) (*domain.Pallet, error) {
	pallet, err := GetSupplyById(product.PalletID)
	if err != nil || pallet == nil {
		log.Print(product.PalletID, 5)
		return nil, err
	}
	log.Print(pallet, 5)
	product.PalletID = pallet.ID
	if err := db.DB.Model(pallet).Association("PalletizedProduct").Append(&product); err != nil {
		return nil, err
	}

	if err := db.DB.Preload("PalletizedProduct").First(pallet, pallet.ID).Error; err != nil {
		return nil, err
	}

	return pallet, nil
}
