// pkg/utils/jwt.go
package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v4"
)

var JwtSecret = []byte("your_very_secret_key_here") // In production, use an environment variable

func GenerateToken(userID, username string) (string, string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours [cite: 983]
	
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtSecret)
	
	return tokenString, expirationTime.Format(time.RFC3339), err
}