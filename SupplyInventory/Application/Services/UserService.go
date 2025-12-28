package services

import (
	errors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"
	userRepo "go_inventory/SupplyInventory/Domain/contracts/repositories/User"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
    CreateUser(name string, email string, password string) (*entities.UserEntity, *errors.AppError)
    Login(email string, password string) (*entities.UserEntity, *errors.AppError)
}

type userService struct {
	Repository userRepo.UserRepository
}

func NewUserService(repository userRepo.UserRepository) UserService {
	return &userService{Repository: repository}
}

func (service *userService) CreateUser(name, email, password string) (*entities.UserEntity, *errors.AppError) {
	emailExists, _ := service.Repository.FindByEmail(email)
	if emailExists != nil {
		return nil, errors.NewAppError("Error checking existing email", 422)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewAppError("Could not hash password", 500)
	}

	user := &entities.UserEntity{
		Name:     name,
		Email:    email,
		Password: string(hashed),
	}

	if err := service.Repository.Create(user); err != nil {
		return nil, errors.NewAppError(err.Error(), 500)
	}

	return user, nil
}

func (service *userService) Login(email string, password string) (*entities.UserEntity, *errors.AppError) {
	user, err := service.Repository.FindByEmail(email)
	if err != nil {
		return nil, errors.NewAppError("user not found", 404)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, errors.NewAppError("invalid credentials", 401)
	}

	return user, nil
}
