package services

import (
	"fmt"
	"os"

	"github.com/skip2/go-qrcode"
)

func CreateQRCode(positionId uint) (string, error) {
	baseUrl := os.Getenv("BASE_URL") + os.Getenv("PORT") + "/positions/%d"
	link := fmt.Sprintf(baseUrl, positionId)

	os.MkdirAll("storage/qrcodes", os.ModePerm)

	qrFile := fmt.Sprintf("storage/qrcodes/position_%d.png", positionId)
	err := qrcode.WriteFile(link, qrcode.Highest, 1024, qrFile)
	if err != nil {
		return "", err
	}

	return qrFile, nil
}
