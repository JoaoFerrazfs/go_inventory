package repositories

import (
	entities "go_inventory/SupplyInventory/Domain/Entities"
)

type UserRepository interface {
	Create(user *entities.UserEntity) error
	FindByEmail(email string) (*entities.UserEntity, error)
}
