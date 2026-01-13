package testutils

import (
	entities "go_inventory/SupplyInventory/Domain/Entities"

	"gorm.io/gorm"
)

// FakeDBAdapter is a minimal test double implementing the DBAdapter interface used by repositories.
type FakeDBAdapter struct {
	CreateFn            func(value interface{}) error
	FirstByIDFn         func(out interface{}, id uint) error
	FindAllFn           func(out interface{}) error
	WhereFirstFn        func(out interface{}, query string, args ...interface{}) error
	DeleteByIDFn        func(model interface{}, id uint) (int64, error)
	SaveFn              func(value interface{}) error
	PreloadFindFn       func(out interface{}, preload string, id ...uint) error
	WherePreloadFindFn  func(out interface{}, preload string, where string, args ...interface{}) error
	GetDBFn             func() *gorm.DB
	AppendAssociationFn func(pallet *entities.PalletEntity, product *entities.PalletizedProductEntity) error
	CountAndPaginatedFindFn func(model interface{}, out interface{}, total *int64, page int, limit int, preloads []string, where string, args ...interface{}) error
}

func (f *FakeDBAdapter) Create(value interface{}) error {
	if f.CreateFn != nil {
		return f.CreateFn(value)
	}
	return nil
}

func (f *FakeDBAdapter) FirstByID(out interface{}, id uint) error {
	if f.FirstByIDFn != nil {
		return f.FirstByIDFn(out, id)
	}
	return nil
}

func (f *FakeDBAdapter) FindAll(out interface{}) error {
	if f.FindAllFn != nil {
		return f.FindAllFn(out)
	}
	return nil
}

func (f *FakeDBAdapter) WhereFirst(out interface{}, query string, args ...interface{}) error {
	if f.WhereFirstFn != nil {
		return f.WhereFirstFn(out, query, args...)
	}
	// No default; tests should set WhereFirstFn when they expect data.
	return nil
}

func (f *FakeDBAdapter) DeleteByID(model interface{}, id uint) (int64, error) {
	if f.DeleteByIDFn != nil {
		return f.DeleteByIDFn(model, id)
	}
	return 1, nil
}

func (f *FakeDBAdapter) Save(value interface{}) error {
	if f.SaveFn != nil {
		return f.SaveFn(value)
	}
	return nil
}

func (f *FakeDBAdapter) PreloadFind(out interface{}, preload string, id ...uint) error {
	if f.PreloadFindFn != nil {
		return f.PreloadFindFn(out, preload, id...)
	}
	return nil
}

func (f *FakeDBAdapter) WherePreloadFind(out interface{}, preload string, where string, args ...interface{}) error {
	if f.WherePreloadFindFn != nil {
		return f.WherePreloadFindFn(out, preload, where, args...)
	}
	return nil
}

func (f *FakeDBAdapter) GetDB() *gorm.DB {
	if f.GetDBFn != nil {
		return f.GetDBFn()
	}
	return nil
}

func (f *FakeDBAdapter) AppendAssociation(pallet *entities.PalletEntity, product *entities.PalletizedProductEntity) error {
	if f.AppendAssociationFn != nil {
		return f.AppendAssociationFn(pallet, product)
	}
	return nil
}

func (f *FakeDBAdapter) CountAndPaginatedFind(model interface{}, out interface{}, total *int64, page int, limit int, preloads []string, where string, args ...interface{}) error {
	if f.CountAndPaginatedFindFn != nil {
		return f.CountAndPaginatedFindFn(model, out, total, page, limit, preloads, where, args...)
	}
	return nil
}
