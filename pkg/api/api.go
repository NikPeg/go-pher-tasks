package api

import (
	"net/http"
)

// Init регистрирует все API-обработчики.
func Init() {
	// Регистрируем обработчик для вычисления следующей даты
	http.HandleFunc("/api/nextdate", nextDateHandler)
}
