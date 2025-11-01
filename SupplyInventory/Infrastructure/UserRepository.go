package infrastructure

import (
	domain "go_inventory/SupplyInventory/Domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.UserEntity) error
	FindByEmail(email string) (*domain.UserEntity, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (repository *userRepository) Create(user *domain.UserEntity) error {
	return repository.db.Create(user).Error
}

func (repository *userRepository) FindByEmail(email string) (*domain.UserEntity, error) {
	var user domain.UserEntity
	if err := repository.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
