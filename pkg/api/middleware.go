package api

import (
	"net/http"
	"os"
)

// authMiddleware проверяет аутентификацию пользователя перед доступом к защищенным эндпоинтам.
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")

		// Если пароль не установлен, аутентификация не требуется.
		if password == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Получаем токен из cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			// Если cookie нет, возвращаем 401 Unauthorized
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// Валидируем JWT
		if !validateJWT(cookie.Value) {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Если все в порядке, передаем управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}
