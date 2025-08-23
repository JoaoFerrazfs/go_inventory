package services

import (
	"log"
	"strings"

	domain "go_inventory/SupplyInventory/Domain"
	infrastructure "go_inventory/SupplyInventory/Infrastructure"
)

func ListPositions() []domain.Position {
	positions, err := infrastructure.GetAllPositions()
	if err != nil {
		return nil
	}

	return positions
}

func FindPositionById(id uint) (*domain.Position, error) {
	position, err := infrastructure.GetSupplyById(id)
	if err != nil {
		return nil, err
	}

	return position, nil
}

func CreatePosition(position domain.Position) (*domain.Position, error) {
	newPosition, err := infrastructure.AddSupply(position)
	if err != nil {
		return nil, err
	}

	qrcodeLink, err := CreateQRCode(newPosition.ID)
	if err != nil {
		return nil, err
	}

	newPosition.QrCode = qrcodeLink
	newPosition.QrCodeUrl = "http://localhost:3000/" + strings.TrimPrefix(qrcodeLink, "storage/")

	positionWithQrCode, err := infrastructure.UpdateSupply(newPosition)
	if err != nil {
		return nil, err
	}
	log.Printf("New position created with QR code: %+v", positionWithQrCode)
	return positionWithQrCode, nil
}

func AddProductsToPositition(product domain.PositionProduct) *domain.Position {
	newPositionProduct, err := infrastructure.AddProductsToPosition(product)
	if err != nil {
		return nil
	}

	return newPositionProduct
}
