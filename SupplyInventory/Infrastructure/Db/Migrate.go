package db

import (
	"log"

	domain "go_inventory/SupplyInventory/Domain"

	"gorm.io/gorm"
)

// Recebe a instância do DB
func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&domain.PalletRackEntity{},
		&domain.PalletEntity{},
		&domain.PalletizedProductEntity{},
	)
	if err != nil {
		log.Fatal("Erro na migration:", err)
	}

	log.Println("Migration realizada com sucesso")
}
