package services

import (
	"os"
	"time"

	errors "go_inventory/Helpers/Errors"
	"go_inventory/SupplyInventory/Domain/constants"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type JWTService interface {
	GenerateToken(userID uint, username string, role constants.UserRole) (string, *errors.AppError)
	GenerateRefreshToken(userID uint, username string, role constants.UserRole) (string, *errors.AppError)
	ValidateToken(tokenString string) (*jwt.Token, *errors.AppError)
	RefreshToken(token string) (string, *errors.AppError)
}

type jwtService struct{}

func NewJWTService() JWTService {
	return &jwtService{}
}

func (service *jwtService) GenerateToken(userID uint, username string, role constants.UserRole) (string, *errors.AppError) {
	return generateTokenWithExpiration(userID, username, role, time.Hour*1, TokenTypeAccess)
}

func (service *jwtService) GenerateRefreshToken(userID uint, username string, role constants.UserRole) (string, *errors.AppError) {
	return generateTokenWithExpiration(userID, username, role, time.Hour*24*7, TokenTypeRefresh)
}

func generateTokenWithExpiration(userID uint, username string, role constants.UserRole, duration time.Duration, tokenType string) (string, *errors.AppError) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"username":  username,
		"role":      role.String(),
		"tokenType": tokenType,
		"exp":       time.Now().Add(duration).Unix(),
	})

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", errors.NewAppError(err.Error(), 500)
	}

	return tokenString, nil
}

func (service *jwtService) ValidateToken(tokenString string) (*jwt.Token, *errors.AppError) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, errors.NewAppError(err.Error(), 401)
	}

	return token, nil
}

func (service *jwtService) RefreshToken(refreshToken string) (string, *errors.AppError) {
	token, appErr := service.ValidateToken(refreshToken)
	if appErr != nil || !token.Valid {
		return "", errors.NewAppError("invalid refresh token", 401)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.NewAppError("invalid token claims", 401)
	}

	userIDFloat, ok := claims["userID"].(float64)
	if !ok {
		return "", errors.NewAppError("invalid or missing userID in token claims", 401)
	}
	userID := uint(userIDFloat)

	username, ok := claims["username"].(string)
	if !ok {
		return "", errors.NewAppError("invalid or missing username in token claims", 401)
	}

	roleStr, ok := claims["role"].(string)
	if !ok {
		return "", errors.NewAppError("invalid or missing role in token claims", 401)
	}
	role := constants.UserRole(roleStr)

	newToken, appErr := service.GenerateToken(userID, username, role)
	if appErr != nil {
		return "", appErr
	}

	return newToken, nil
}
