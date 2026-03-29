//go:build unit

package services_test

import (
	"os"
	"path/filepath"
	"testing"

	qrCodeService "go_inventory/SupplyInventory/Application/Services/QrCode"
	testutils "go_inventory/SupplyInventory/tests/testutils"

	"github.com/stretchr/testify/assert"
)

func TestCreateQRCode_Success(t *testing.T) {
	// Set
	restoreBase := testutils.SetEnvAndRestore("BASE_URL", "http://localhost:")
	defer restoreBase()
	restorePort := testutils.SetEnvAndRestore("PORT", "3000")
	defer restorePort()

	tempDir, cleanup := testutils.CreateTempDir(t, "qrcode_test")
	defer cleanup()

	// ensure we call service from repo root so DIRECTORY resolves relative to CWD
	old, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer os.Chdir(old)

	svc := qrCodeService.NewQRCodeService()
	path, _, err := svc.CreateQRCode(7)

	// Assertions
	assert.Nil(t, err)
	assert.NotEmpty(t, path)
	// file exists
	_, statErr := os.Stat(filepath.Clean(path))
	assert.Nil(t, statErr)
}

func TestCreateQRCode_DirectoryError(t *testing.T) {
	// Set
	restoreBase := testutils.SetEnvAndRestore("BASE_URL", "http://localhost:")
	defer restoreBase()
	restorePort := testutils.SetEnvAndRestore("PORT", "3000")
	defer restorePort()

	// Make DIRECTORY path unwritable by creating a file with that name
	_ = os.RemoveAll("storage")
	f, _ := os.Create("storage")
	f.Close()
	defer os.Remove("storage")

	svc := qrCodeService.NewQRCodeService()
	_, _, err := svc.CreateQRCode(8)

	// Assertions
	assert.NotNil(t, err)
}
