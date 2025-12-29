package db

import (
	entities "go_inventory/SupplyInventory/Domain/Entities"

	"gorm.io/gorm"
)

// DBAdapter defines the minimal DB operations used by repositories.
type DBAdapter interface {
    Create(value interface{}) error
    FirstByID(out interface{}, id uint) error
    FindAll(out interface{}) error
    WhereFirst(out interface{}, query string, args ...interface{}) error
    DeleteByID(model interface{}, id uint) (int64, error)
    Save(value interface{}) error
    PreloadFind(out interface{}, preload string, id ...uint) error
    AppendAssociation(pallet *entities.PalletEntity, product *entities.PalletizedProductEntity) error
}

// gormAdapter implements DBAdapter using *gorm.DB
type gormAdapter struct{
    db *gorm.DB
}

// NewGormAdapter wraps a *gorm.DB
func NewGormAdapter(db *gorm.DB) DBAdapter {
    return &gormAdapter{db: db}
}

func (g *gormAdapter) Create(value interface{}) error {
    return g.db.Create(value).Error
}

func (g *gormAdapter) FirstByID(out interface{}, id uint) error {
    return g.db.First(out, id).Error
}

func (g *gormAdapter) FindAll(out interface{}) error {
    return g.db.Find(out).Error
}

func (g *gormAdapter) WhereFirst(out interface{}, query string, args ...interface{}) error {
    return g.db.Where(query, args...).First(out).Error
}

func (g *gormAdapter) DeleteByID(model interface{}, id uint) (int64, error) {
    res := g.db.Delete(model, id)
    return res.RowsAffected, res.Error
}

func (g *gormAdapter) Save(value interface{}) error {
    return g.db.Save(value).Error
}

func (g *gormAdapter) PreloadFind(out interface{}, preload string, id ...uint) error {
    if len(id) > 0 {
        return g.db.Preload(preload).First(out, id[0]).Error
    }
    return g.db.Preload(preload).Find(out).Error
}

func (g *gormAdapter) AppendAssociation(pallet *entities.PalletEntity, product *entities.PalletizedProductEntity) error {
    return g.db.Model(pallet).Association("PalletizedProduct").Append(product)
}
