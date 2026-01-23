package entities

type ProductEntity struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	EAN  string `gorm:"unique;not null;size:13" json:"ean" binding:"required"`
	Name string `gorm:"not null" json:"name" binding:"required"`
}

func (ProductEntity) TableName() string {
	return "products"
}

func NewProductEntity(ean string, name string) *ProductEntity {
	return &ProductEntity{
		EAN:  ean,
		Name: name,
	}
}
