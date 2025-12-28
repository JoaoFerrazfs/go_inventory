package testutils

import (
	entities "go_inventory/SupplyInventory/Domain/Entities"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateTestUser(db *gorm.DB, name, email, password string) *entities.UserEntity {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := entities.NewUserEntity(name, email, string(hashed))
	db.Create(user)
	return user
}

func CreateTestPallet(db *gorm.DB, name string, palletRackID uint) *entities.PalletEntity {
	pallet := &entities.PalletEntity{
		Name:         name,
		PalletRackID: palletRackID,
	}
	db.Create(pallet)
	return pallet
}

func CreateTestPalletRack(db *gorm.DB, name, location string, totalCapacity int) *entities.PalletRackEntity {
	rack := &entities.PalletRackEntity{
		Name:          name,
		Location:      location,
		TotalCapacity: totalCapacity,
	}
	db.Create(rack)
	return rack
}

// Add more fixtures as needed for other entities