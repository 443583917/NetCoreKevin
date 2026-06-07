package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	secret := "test-secret-key"
	expireHour := 24

	token, err := GenerateToken(secret, expireHour, 1, "testuser", 1000)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestParseToken(t *testing.T) {
	secret := "test-secret-key"
	expireHour := 24

	token, err := GenerateToken(secret, expireHour, 1, "testuser", 1000)
	assert.NoError(t, err)

	claims, err := ParseToken(secret, token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, int64(1), claims.UserID)
	assert.Equal(t, "testuser", claims.UserName)
	assert.Equal(t, int64(1000), claims.TenantID)
}

func TestParseToken_Expired(t *testing.T) {
	secret := "test-secret-key"

	claims := &JWTClaims{
		UserID:   1,
		UserName: "testuser",
		TenantID: 1000,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	parsedClaims, err := ParseToken(secret, tokenString)
	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
}

func TestParseToken_InvalidSecret(t *testing.T) {
	secret := "test-secret-key"
	wrongSecret := "wrong-secret-key"

	token, err := GenerateToken(secret, 24, 1, "testuser", 1000)
	assert.NoError(t, err)

	claims, err := ParseToken(wrongSecret, token)
	assert.Error(t, err)
	assert.Nil(t, claims)
}
