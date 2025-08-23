package services

import (
	"fmt"
	"os"

	"github.com/skip2/go-qrcode"
)

// Interface do QRCodeService
type QRCodeService interface {
	CreateQRCode(palletId uint) (string, error)
}

// Implementação
type qrCodeService struct{}

// Construtor
func NewQRCodeService() QRCodeService {
	return &qrCodeService{}
}

// Método
func (s *qrCodeService) CreateQRCode(palletId uint) (string, error) {
	baseUrl := os.Getenv("BASE_URL") + os.Getenv("PORT") + "/pallets/%d"
	link := fmt.Sprintf(baseUrl, palletId)

	// Cria pasta caso não exista
	if err := os.MkdirAll("storage/qrcodes", os.ModePerm); err != nil {
		return "", err
	}

	qrFile := fmt.Sprintf("storage/qrcodes/pallet_%d.png", palletId)
	err := qrcode.WriteFile(link, qrcode.Highest, 1024, qrFile)
	if err != nil {
		return "", err
	}

	return qrFile, nil
}
