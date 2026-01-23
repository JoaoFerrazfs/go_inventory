package db

import (
	entities "go_inventory/SupplyInventory/Domain/Entities"

	"gorm.io/gorm"
)

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
