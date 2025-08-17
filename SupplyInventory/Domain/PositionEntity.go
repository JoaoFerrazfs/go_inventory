package domain

type Position struct {
	ID       uint              `gorm:"primaryKey" json:"id"`
	Name     string            `json:"name" binding:"required"`
	Products []PositionProduct `gorm:"foreignKey:PositionID" json:"products"`
}
