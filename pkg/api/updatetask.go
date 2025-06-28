package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go1f/pkg/db"
	"net/http"
	"strconv"
)

// UpdateTaskRequest - это структура для разбора входящего JSON при обновлении.
// ID здесь может быть строкой, чтобы соответствовать тому, что присылает фронтенд.
type UpdateTaskRequest struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateTaskRequest
	// 1. Декодируем запрос в нашу временную структуру UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	// 2. Валидируем и конвертируем ID из строки в число
	if req.ID == "" {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "id is required for update"})
		return
	}
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid id format"})
		return
	}

	// 3. Собираем "настоящую" структуру db.Task для работы с БД
	task := db.Task{
		ID:      id,
		Date:    req.Date,
		Title:   req.Title,
        Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	// 4. Используем ту же функцию валидации, что и раньше
	if err := validateAndProcessTask(&task); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// 5. Обновляем задачу в БД
	if err := db.UpdateTask(task); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONResponse(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		} else {
			writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to update task"})
		}
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{})
}
