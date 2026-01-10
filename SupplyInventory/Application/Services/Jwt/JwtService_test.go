package services_test

import (
	"testing"
	"time"

	testutils "go_inventory/SupplyInventory/tests/testutils"

	jwtService "go_inventory/SupplyInventory/Application/Services/Jwt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken_Success(t *testing.T) {
	// Set
	restore := testutils.SetEnvAndRestore("JWT_SECRET", "testsecret")
	defer restore()

	// Actions
	svc := jwtService.NewJWTService()
	tokenStr, appErr := svc.GenerateToken(10, "john")

	// Assertions
	assert.Nil(t, appErr)
	assert.NotEmpty(t, tokenStr)

	// parse to ensure valid signature
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("testsecret"), nil
	})
	assert.Nil(t, err)
	assert.True(t, token.Valid)
}

func TestGenerateRefreshToken_Success(t *testing.T) {
	// Set
	restore := testutils.SetEnvAndRestore("JWT_SECRET", "testsecret")
	defer restore()

	// Actions
	svc := jwtService.NewJWTService()
	tokenStr, appErr := svc.GenerateRefreshToken(20, "mary")

	// Assertions
	assert.Nil(t, appErr)
	assert.NotEmpty(t, tokenStr)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("testsecret"), nil
	})
	assert.Nil(t, err)
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		assert.Equal(t, "refresh", claims["tokenType"])
	} else {
		t.Fatalf("invalid claims type")
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	// Set
	restore := testutils.SetEnvAndRestore("JWT_SECRET", "testsecret")
	defer restore()

	svc := jwtService.NewJWTService()

	// Actions
	_, appErr := svc.ValidateToken("this.is.not.a.token")

	// Assertions
	assert.NotNil(t, appErr)
}

func TestValidateToken_ValidToken(t *testing.T) {
	// Set
	restore := testutils.SetEnvAndRestore("JWT_SECRET", "testsecret")
	defer restore()

	// create valid token manually
	claims := jwt.MapClaims{
		"userID":    5,
		"username":  "alice",
		"tokenType": "access",
		"exp":       time.Now().Add(time.Minute * 5).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("testsecret"))

	svc := jwtService.NewJWTService()

	// Actions
	parsed, appErr := svc.ValidateToken(tokenStr)

	// Assertions
	assert.Nil(t, appErr)
	assert.NotNil(t, parsed)
	assert.True(t, parsed.Valid)
}

func TestRefreshToken_Success(t *testing.T) {
	// Set
	restore := testutils.SetEnvAndRestore("JWT_SECRET", "testsecret")
	defer restore()

	// create valid refresh token
	claims := jwt.MapClaims{
		"userID":    9,
		"username":  "bob",
		"tokenType": "refresh",
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("testsecret"))

	svc := jwtService.NewJWTService()

	// Actions
	newToken, appErr := svc.RefreshToken(tokenStr)

	// Assertions
	assert.Nil(t, appErr)
	assert.NotEmpty(t, newToken)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	// Set
	restore := testutils.SetEnvAndRestore("JWT_SECRET", "testsecret")
	defer restore()

	svc := jwtService.NewJWTService()

	// Actions
	_, appErr := svc.RefreshToken("invalid.token")

	// Assertions
	assert.NotNil(t, appErr)
}
