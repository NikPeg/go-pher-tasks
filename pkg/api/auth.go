package api

import (
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SecretKey - ключ для подписи токена.
var SecretKey = []byte("super_secret_key_change_it_in_production")

// Claims определяет структуру полезной нагрузки (payload) нашего JWT.
type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

// generateJWT создает новый JWT для заданного пароля.
func generateJWT(password string) (string, error) {
	// Создаем хэш пароля, чтобы не хранить его в открытом виде в токене.
	h := sha256.New()
	h.Write([]byte(password))
	passwordHash := fmt.Sprintf("%x", h.Sum(nil))

	// Устанавливаем время жизни токена (например, 8 часов)
	expirationTime := time.Now().Add(8 * time.Hour)

	// Создаем полезную нагрузку (claims)
	claims := &Claims{
		PasswordHash: passwordHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Создаем новый токен с нашими claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен нашим секретным ключом
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// validateJWT проверяет JWT на валидность и соответствие текущему паролю.
func validateJWT(tokenString string) bool {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Убеждаемся, что используется правильный алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return SecretKey, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	// Самая важная проверка: сверяем хэш пароля из токена
	// с хэшем текущего пароля из переменной окружения.
	// Это делает токен невалидным, если пароль был изменен.
	currentPassword := os.Getenv("TODO_PASSWORD")
	h := sha256.New()
	h.Write([]byte(currentPassword))
	currentPasswordHash := fmt.Sprintf("%x", h.Sum(nil))

	return claims.PasswordHash == currentPasswordHash
}
