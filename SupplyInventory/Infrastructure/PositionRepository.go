package infrastructure

import (
	"log"

	domain "go_inventory/SupplyInventory/Domain"
	db "go_inventory/SupplyInventory/Infrastructure/Db"
)

func GetAllPositions() ([]domain.Position, error) {
	var positions []domain.Position
	if err := db.DB.Preload("Products").Find(&positions).Error; err != nil {
		return nil, err
	}
	return positions, nil
}

func GetSupplyById(id uint) (*domain.Position, error) {
	var position domain.Position
	if err := db.DB.Preload("Products").First(&position, id).Error; err != nil {
		return nil, err
	}
	log.Printf("Position found: %+v", position)
	return &position, nil
}

func AddSupply(position domain.Position) (*domain.Position, error) {
	if err := db.DB.Create(&position).Error; err != nil {
		return nil, err
	}
	return &position, nil
}

func AddProductsToPosition(product domain.PositionProduct) (*domain.Position, error) {
	// 1️⃣ Busca a posição com os produtos já existentes
	position, err := GetSupplyById(product.PositionID)
	if err != nil || position == nil {
		return nil, err
	}

	// 2️⃣ Garante que o PositionID do produto esteja correto
	product.PositionID = position.ID

	// 3️⃣ Adiciona o produto usando a associação
	if err := db.DB.Model(position).Association("Products").Append(&product); err != nil {
		return nil, err
	}

	// 4️⃣ Recarrega os produtos da posição
	if err := db.DB.Preload("Products").First(position, position.ID).Error; err != nil {
		return nil, err
	}

	return position, nil
}

func UpdateSupply(position domain.Position) (*domain.Position, error) {
	if err := db.DB.Save(&position).Error; err != nil {
		return nil, err
	}
	return &position, nil
}
