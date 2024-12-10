package services

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateToken генерирует JWT для пользователя.
func GenerateToken(userID int) (string, error) {
	// Получаем секретный ключ из переменных окружения.
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY is not set in the environment variables")
	}

	// Создаем токен с пользовательскими данными.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Истечение через 24 часа.
	})

	// Подписываем токен секретным ключом.
	return token.SignedString([]byte(secretKey))
}
