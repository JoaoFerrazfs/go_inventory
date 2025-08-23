package services

import (
	"fmt"
	"os"

	"github.com/skip2/go-qrcode"
)

func CreateQRCode(palletId uint) (string, error) {
	baseUrl := os.Getenv("BASE_URL") + os.Getenv("PORT") + "/pallets/%d"
	link := fmt.Sprintf(baseUrl, palletId)

	os.MkdirAll("storage/qrcodes", os.ModePerm)

	qrFile := fmt.Sprintf("storage/qrcodes/pallet_%d.png", palletId)
	err := qrcode.WriteFile(link, qrcode.Highest, 1024, qrFile)
	if err != nil {
		return "", err
	}

	return qrFile, nil
}
