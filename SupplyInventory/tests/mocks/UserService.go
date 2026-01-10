package mocks

import (
	appErrors "go_inventory/Helpers/Errors"
	entities "go_inventory/SupplyInventory/Domain/Entities"

	"github.com/stretchr/testify/mock"
)

type UserService struct {
	mock.Mock
}

func (m *UserService) CreateUser(name string, email string, password string) (*entities.UserEntity, *appErrors.AppError) {
	args := m.Called(name, email, password)
	result := args.Get(0)
	if result == nil {
		return nil, args.Get(1).(*appErrors.AppError)
	}
	return result.(*entities.UserEntity), args.Get(1).(*appErrors.AppError)
}

func (m *UserService) Login(email string, password string) (*entities.UserEntity, *appErrors.AppError) {
	args := m.Called(email, password)
	result := args.Get(0)
	if result == nil {
		return nil, args.Get(1).(*appErrors.AppError)
	}
	return result.(*entities.UserEntity), args.Get(1).(*appErrors.AppError)
}

func (m *UserService) GetUserByID(id uint) (*entities.UserEntity, *appErrors.AppError) {
	args := m.Called(id)
	result := args.Get(0)
	if result == nil {
		return nil, args.Get(1).(*appErrors.AppError)
	}
	return result.(*entities.UserEntity), args.Get(1).(*appErrors.AppError)
}