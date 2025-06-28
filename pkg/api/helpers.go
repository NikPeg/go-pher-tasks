package api

import (
	"encoding/json"
	"net/http"
)

// writeJSONResponse отправляет ответ в формате JSON с указанным статус-кодом.
func writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// В случае ошибки кодирования, логируем и отправляем простую ошибку
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
	}
}
