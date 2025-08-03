package infrastructure

import (
	"go_inventory/SupplyInventory/Domain"
    "go_inventory/SupplyInventory/Infrastructure/Db"
)


func GetAllPositions() ([]domain.Position, error) {
    rows, err := db.DB.Query("SELECT id, name, stock FROM positions")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var positions []domain.Position
    for rows.Next() {
        var position domain.Position
        if err := rows.Scan(&position.ID, &position.Name, &position.Stock); err != nil {
            return nil, err
        }
        positions = append(positions, position)
    }
    return positions, nil
}

func GetSupplyById(id string) (*domain.Position, error) {
    row := db.DB.QueryRow("SELECT id, name, stock FROM positions WHERE id = ?", id)
    var position domain.Position

    if err := row.Scan(&position.ID, &position.Name, &position.Stock); err != nil {
        return nil, err
    }
    return &position, nil
}

func AddSupply(position domain.Position) (*domain.Position, error) {
    result, err := db.DB.Exec(
        "INSERT INTO positions (name, stock, ean) VALUES (?, ?, ?)",
        position.Name, position.Stock, position.EAN ,
    )
    if err != nil {
        return nil, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return nil, err
    }

    newPosition := &domain.Position{
        ID:    int(id),
        Name:  position.Name,
        Stock: position.Stock,
		EAN:   int(position.EAN),
    }
    return newPosition, nil
}