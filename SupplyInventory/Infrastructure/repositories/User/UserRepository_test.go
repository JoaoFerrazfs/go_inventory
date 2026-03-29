//go:build unit

package infrastructure_test

import (
	"testing"

	entities "go_inventory/SupplyInventory/Domain/Entities"
	repos "go_inventory/SupplyInventory/Infrastructure/repositories/User"
	"go_inventory/SupplyInventory/tests/testutils"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	repo := repos.NewUserRepository(adapter)
	u := &entities.UserEntity{Name: "Unit", Email: "unit@example.com", Password: "pwd"}

	// Actions
	err := repo.Create(u)

	// Assertions
	assert.Nil(t, err)
}

func TestUserRepository_FindByEmail_Unit(t *testing.T) {
	// Set
	adapter := &testutils.FakeDBAdapter{}
	adapter.WhereFirstFn = func(out interface{}, query string, args ...interface{}) error {
		if u, ok := out.(*entities.UserEntity); ok {
			u.ID = 1
			u.Name = "TestUser"
			if len(args) > 0 {
				if email, ok := args[0].(string); ok {
					u.Email = email
				}
			}
		}
		return nil
	}
	repo := repos.NewUserRepository(adapter)

	// Actions
	got, err := repo.FindByEmail("test@example.com")

	// Assertions
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, uint(1), got.ID)
	assert.Equal(t, "test@example.com", got.Email)
}
