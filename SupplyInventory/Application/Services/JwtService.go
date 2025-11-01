package services

import (
	"os"
	"time"

	errors "go_inventory/Helpers/Errors"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userID uint, username string) (string, *errors.AppError)
	ValidateToken(tokenString string) (bool, *errors.AppError)
}

type jwtService struct{}

func NewJWTService() JWTService {
	return &jwtService{}
}

func (service *jwtService) GenerateToken(userID uint, username string) (string, *errors.AppError) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // expira em 1h
	})

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", errors.NewAppError(err.Error(), 500)
	}

	return tokenString, nil
}

func (service *jwtService) ValidateToken(tokenString string) (bool, *errors.AppError) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return jwtSecret, nil
	})
	if err != nil {
		return false, errors.NewAppError(err.Error(), 401)
	}

	return token.Valid, nil
}
