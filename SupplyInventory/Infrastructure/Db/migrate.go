package db

import (
	"log"

	domain "go_inventory/SupplyInventory/Domain"
)

func Migrate() {
	err := DB.AutoMigrate(
		&domain.Pallet{},
		&domain.PalletizedProductEntity{},
	)
	if err != nil {
		log.Fatal("Erro na migration:", err)
	}

	log.Println("Migration realizada com sucesso")
}
