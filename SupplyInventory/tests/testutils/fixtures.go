package testutils

import (
	entities "go_inventory/SupplyInventory/Domain/Entities"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateTestUser(db *gorm.DB, name, email, password string) *entities.UserEntity {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Check if user already exists
	var existingUser entities.UserEntity
	if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return &existingUser // Return existing user if found
	}

	user := entities.NewUserEntity(name, email, string(hashed))
	if err := db.Create(user).Error; err != nil {
		panic(err) // Ensure the user creation error is handled
	}
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
	inv := CreateTestInventory(db)
	if inv == nil {
		panic("Failed to create test inventory")
	}

	rack := &entities.PalletRackEntity{
		InventoryID:   inv.ID,
		Name:          name,
		Location:      location,
		TotalCapacity: totalCapacity,
	}
	if err := db.Create(rack).Error; err != nil {
		panic(err)
	}
	return rack
}

func CreateTestInventory(db *gorm.DB) *entities.InventoryEntity {
	user := CreateTestUser(db, "Test User", "test+inventory@example.com", "password")
	if user == nil {
		panic("Failed to create test user")
	}

	inv := &entities.InventoryEntity{
		Name:      "Test Inventory",
		UserID:    user.ID,
		Status:    "Open",
		StartedAt: time.Now(),
	}
	if err := db.Create(inv).Error; err != nil {
		panic(err) // Ensure the inventory creation error is handled
	}
	return inv
}

// Add more fixtures as needed for other entities
