package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oojoseph67/ecommerce/internal/config"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Generates access and refresh tokens
func GenerateTokenPair(config *config.JWTConfig, userId uint, email, role string) (accessToken, refreshToken string, err error) {
	// creating accessToken
	accessClaims := &Claims{
		UserID: userId,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodES256, accessClaims)
	accessTokenString, err := at.SignedString([]byte(config.Secret))

	if err != nil {
		return "", "", err
	}

	// creating refreshToken
	refreshClaims := &Claims{
		UserID: userId,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.RefreshTokenExpires)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodES256, refreshClaims)
	refreshTokenString, err := rt.SignedString([]byte(config.Secret))

	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// Checks if accessToken or refreshToken is valid
func ValidateToken(tokenString, secret string) (claims *Claims, err error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
