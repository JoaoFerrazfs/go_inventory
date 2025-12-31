package services

import (
	"fmt"
	"os"
	"path/filepath"

	storage "go_inventory/SupplyInventory/Application/Services/Storage"

	"github.com/skip2/go-qrcode"
)

type QRCodeService interface {
	// returns storage path and public URL
	CreateQRCode(palletId uint) (string, string, error)
}

type qrCodeService struct{
	storage storage.Storage
	qrcodeDir string // relative path inside storage, e.g. "qrcodes"
}

const DEFAULT_DIRECTORY = "storage/qrcodes"

func NewQRCodeService() QRCodeService {
	// default to local storage using filesystem
	baseURL := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")
	if baseURL == "" {
		baseURL = "http://localhost:3000/"
	}
	if port != "" && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}
	// Use local storage pointed at ./storage
	local := storage.NewLocalStorage("storage", baseURL+"/")
	return &qrCodeService{storage: local, qrcodeDir: "qrcodes"}
}

func NewQRCodeServiceWithStorage(st storage.Storage) QRCodeService{
	return &qrCodeService{storage: st, qrcodeDir: "qrcodes"}
}

func (service *qrCodeService) CreateQRCode(palletId uint) (string, string, error) {
	baseUrl := os.Getenv("BASE_URL") + os.Getenv("PORT") + "/pallets/%d"
	link := fmt.Sprintf(baseUrl, palletId)

	// generate png into a buffer
	filename := fmt.Sprintf("pallet_%d.png", palletId)
	relPath := filepath.Join(service.qrcodeDir, filename)

	// create temp file
	tmpFile := filepath.Join(os.TempDir(), filename)
	if err := qrcode.WriteFile(link, qrcode.Highest, 1024, tmpFile); err != nil {
		return "", "", err
	}
	defer os.Remove(tmpFile)

	f, err := os.Open(tmpFile)
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	_, err = service.storage.Upload(relPath, f)
	if err != nil {
		return "", "", err
	}

	publicURL, _ := service.storage.GetURL(relPath)

	// For compatibility with previous behavior, if using local storage return the full fs path
	if ls, ok := service.storage.(*storage.LocalStorage); ok {
		full := filepath.Join(ls.BaseDir, relPath)
		return full, publicURL, nil
	}

	return relPath, publicURL, nil
}
