package api

import (
	"encoding/json"
	"net/http"
	"os"
)

type signinRequest struct {
	Password string `json:"password"`
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" {
		writeJSONResponse(w, http.StatusServiceUnavailable, map[string]string{"error": "authentication is not enabled"})
		return
	}

	var req signinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	if req.Password != expectedPassword {
		writeJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "invalid password"})
		return
	}

	// Пароль верный, генерируем токен
	token, err := generateJWT(expectedPassword)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"token": token})
}
