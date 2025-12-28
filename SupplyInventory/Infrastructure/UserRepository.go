package infrastructure

import (
	entities "go_inventory/SupplyInventory/Domain/Entities"
	repositories "go_inventory/SupplyInventory/Domain/contracts/repositories/User"

	"gorm.io/gorm"
)


type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{db: db}
}

func (repository *userRepository) Create(user *entities.UserEntity) error {
	return repository.db.Create(user).Error
}

func (repository *userRepository) FindByEmail(email string) (*entities.UserEntity, error) {
       var user entities.UserEntity
       if err := repository.db.Where("email = ?", email).First(&user).Error; err != nil {
	       return nil, err
       }
       return &user, nil
}
