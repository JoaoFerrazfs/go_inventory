package db

import (
	entities "go_inventory/SupplyInventory/Domain/Entities"

	"gorm.io/gorm"
)

// Migrate runs all database migrations
// Note: With RBAC implementation, UserEntity now includes Role field
// If you have existing database, manually add column:
// ALTER TABLE user_entities ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user' AFTER password;
func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entities.InventoryEntity{},
		&entities.PalletRackEntity{},
		&entities.PalletEntity{},
		&entities.PalletizedProductEntity{},
		&entities.UserEntity{},
		&entities.ProductEntity{},
	)
	if err != nil {
		return err
	}

	return nil
}
