package services

import (
	"testing"

	entities "go_inventory/SupplyInventory/Domain/Entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) Create(user *entities.UserEntity) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockUserRepo) FindByEmail(email string) (*entities.UserEntity, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.UserEntity), args.Error(1)
}

func TestCreateUser_Success(t *testing.T) {
	// Set
	repo := &mockUserRepo{}
	svc := NewUserService(repo)

	// Expectations
	repo.On("FindByEmail", "new@example.com").Return(nil, assert.AnError)
	repo.On("Create", mock.Anything).Return(nil)

	// Actions
	user, err := svc.CreateUser("New", "new@example.com", "password123")

	// Assertions
	assert.NotNil(t, user)
	assert.Nil(t, err)
	repo.AssertExpectations(t)
}

func TestCreateUser_EmailExists(t *testing.T) {
	// Set
	repo := &mockUserRepo{}
	svc := NewUserService(repo)

	// Expectations: FindByEmail returns user (exists)
	existing := &entities.UserEntity{ID: 1, Email: "exists@example.com"}
	repo.On("FindByEmail", "exists@example.com").Return(existing, nil)

	// Actions
	user, err := svc.CreateUser("Exist", "exists@example.com", "pass")

	// Assertions
	assert.Nil(t, user)
	assert.NotNil(t, err)
	repo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	// Set
	repo := &mockUserRepo{}
	svc := NewUserService(repo)

	// Create a hashed password using bcrypt to match login expectations
	hashed, _ := bcrypt.GenerateFromPassword([]byte("secretpass"), bcrypt.DefaultCost)
	stored := &entities.UserEntity{ID: 2, Email: "login@example.com", Password: string(hashed)}
	repo.On("FindByEmail", "login@example.com").Return(stored, nil)

	// Actions
	user, err := svc.Login("login@example.com", "secretpass")

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, stored, user)
	repo.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	// Set
	repo := &mockUserRepo{}
	svc := NewUserService(repo)

	stored := &entities.UserEntity{ID: 3, Email: "bad@example.com", Password: "notahash"}
	repo.On("FindByEmail", "bad@example.com").Return(stored, nil)

	// Actions
	user, err := svc.Login("bad@example.com", "wrongpass")

	// Assertions
	assert.Nil(t, user)
	assert.NotNil(t, err)
	repo.AssertExpectations(t)
}
