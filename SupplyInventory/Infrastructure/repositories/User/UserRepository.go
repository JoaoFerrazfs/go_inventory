package infrastructure

import (
	entities "go_inventory/SupplyInventory/Domain/Entities"
	repositories "go_inventory/SupplyInventory/Domain/contracts/repositories/User"
	dbadapter "go_inventory/SupplyInventory/Infrastructure/repositories/db"
)


type userRepository struct {
	db dbadapter.DBAdapter
}

func NewUserRepository(db dbadapter.DBAdapter) repositories.UserRepository {
	return &userRepository{db: db}
}

func (repository *userRepository) Create(user *entities.UserEntity) error {
	return repository.db.Create(user)
}

func (repository *userRepository) FindByEmail(email string) (*entities.UserEntity, error) {
	   var user entities.UserEntity
	   if err := repository.db.WhereFirst(&user, "email = ?", email); err != nil {
			   return nil, err
	   }
	   return &user, nil
}

func (repository *userRepository) FindByID(id uint) (*entities.UserEntity, error) {
	var user entities.UserEntity
	if err := repository.db.FirstByID(&user, id); err != nil {
		return nil, err
	}
	return &user, nil
}