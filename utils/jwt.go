package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = []byte("supersecretjwt") // sama dengan .env

func GenerateJWT(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":   userID,
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWT_SECRET)
}
