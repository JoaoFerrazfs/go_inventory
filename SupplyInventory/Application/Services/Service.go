package services

import (
	"log"
	"strings"

	domain "go_inventory/SupplyInventory/Domain"
	infrastructure "go_inventory/SupplyInventory/Infrastructure"
)

func ListPallets() []domain.Pallet {
	pallets, err := infrastructure.GetAllPallets()
	if err != nil {
		return nil
	}

	return pallets
}

func FindPalletById(id uint) (*domain.Pallet, error) {
	pallet, err := infrastructure.GetSupplyById(id)
	if err != nil {
		return nil, err
	}

	return pallet, nil
}

func CreatePallet(pallet domain.Pallet) (*domain.Pallet, error) {
	newPallet, err := infrastructure.AddSupply(pallet)
	if err != nil {
		return nil, err
	}

	qrcodeLink, err := CreateQRCode(newPallet.ID)
	if err != nil {
		return nil, err
	}

	newPallet.QrCode = qrcodeLink
	newPallet.QrCodeUrl = "http://localhost:3000/" + strings.TrimPrefix(qrcodeLink, "storage/")

	palletWithQrCode, err := infrastructure.UpdateSupply(newPallet)
	if err != nil {
		return nil, err
	}
	log.Printf("New pallet created with QR code: %+v", palletWithQrCode)
	return palletWithQrCode, nil
}

func AddProductsToPallet(product domain.PalletizedProductEntity) *domain.Pallet {
	newPalletProduct, err := infrastructure.AddProductsToPallet(product)
	if err != nil {
		return nil
	}

	return newPalletProduct
}
