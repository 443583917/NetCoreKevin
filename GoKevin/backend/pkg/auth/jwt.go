package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   int64  `json:"userId"`
	UserName string `json:"userName"`
	TenantID int64  `json:"tenantId"`
	jwt.RegisteredClaims
}

func GenerateToken(secret string, expireHour int, userID int64, userName string, tenantID int64) (string, error) {
	claims := &JWTClaims{
		UserID:   userID,
		UserName: userName,
		TenantID: tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-kevin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(secret string, tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func RefreshToken(secret string, expireHour int, tokenString string) (string, error) {
	claims, err := ParseToken(secret, tokenString)
	if err != nil {
		return "", err
	}

	if claims.ExpiresAt != nil {
		remaining := time.Until(claims.ExpiresAt.Time)
		total := time.Duration(expireHour) * time.Hour
		if remaining > total/3 {
			return "", errors.New("token not eligible for refresh")
		}
	}

	return GenerateToken(secret, expireHour, claims.UserID, claims.UserName, claims.TenantID)
}
