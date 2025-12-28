package entities

import "time"

type UserEntity struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Email     string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserEntity(
	name string,
	email string,
	password string,
) *UserEntity {
	return &UserEntity{
		Name:     name,
		Email:    email,
		Password: password,
	}
}
