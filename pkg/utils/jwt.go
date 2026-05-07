package utils

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// genereate token
func GenerateToken(userID int, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("CRITICAL: JWT_SECRET tidak ditemukan di environment")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// verify token
func VerifyToken(tokenString string) (int, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return 0, errors.New("CRITICAL: JWT_SECRET tidak ditemukan di environment")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode penandatanganan token tidak valid")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("token tidak valid atau sudah kedaluwarsa")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("gagal membaca klaim token")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("payload token tidak memiliki user_id")
	}

	return int(userIDFloat), nil
}
