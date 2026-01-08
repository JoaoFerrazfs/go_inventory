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
	// derive inventory id from pallet rack
	var rack entities.PalletRackEntity
	result := db.First(&rack, palletRackID)
	if result.Error != nil {
		panic(result.Error)
	}
	pallet := &entities.PalletEntity{
		Name:         name,
		PalletRackID: palletRackID,
		InventoryID:  rack.InventoryID,
	}
	db.Create(pallet)
	return pallet
}

func CreateTestPalletRack(db *gorm.DB, name, location string, totalCapacity int) *entities.PalletRackEntity {
	// ensure an inventory exists for the rack
	inv := CreateTestInventory(db)
	rack := &entities.PalletRackEntity{
		InventoryID:   inv.ID,
		Name:          name,
		Location:      location,
		TotalCapacity: totalCapacity,
	}
	db.Create(rack)
	return rack
}

func CreateTestInventory(db *gorm.DB) *entities.InventoryEntity {
	inv := &entities.InventoryEntity{}
	db.Create(inv)
	return inv
}

// Add more fixtures as needed for other entities
