package services

import (
	"fmt"
	"os"

	"github.com/skip2/go-qrcode"
)

type QRCodeService interface {
	CreateQRCode(palletId uint) (string, error)
}

type qrCodeService struct{}

const DIRECTORY = "storage/qrcodes"

func NewQRCodeService() QRCodeService {
	return &qrCodeService{}
}

func (service *qrCodeService) CreateQRCode(palletId uint) (string, error) {
	baseUrl := os.Getenv("BASE_URL") + os.Getenv("PORT") + "/pallets/%d"
	link := fmt.Sprintf(baseUrl, palletId)

	if err := os.MkdirAll(DIRECTORY, os.ModePerm); err != nil {
		return "", err
	}

	qrFile := fmt.Sprintf(DIRECTORY+"/pallet_%d.png", palletId)
	err := qrcode.WriteFile(link, qrcode.Highest, 1024, qrFile)
	if err != nil {
		return "", err
	}

	return qrFile, nil
}
