package mocks

import (
	entities "go_inventory/SupplyInventory/Domain/Entities"

	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Create(user *entities.UserEntity) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepository) FindByEmail(email string) (*entities.UserEntity, error) {
	args := m.Called(email)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*entities.UserEntity), args.Error(1)
}

func (m *UserRepository) FindByID(id uint) (*entities.UserEntity, error) {
	args := m.Called(id)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*entities.UserEntity), args.Error(1)
}