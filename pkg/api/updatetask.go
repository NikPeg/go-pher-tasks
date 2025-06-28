package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go1f/pkg/db"
	"net/http"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	if task.ID == 0 {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "id is required for update"})
		return
	}

	// Используем ту же функцию валидации, что и при добавлении задачи
	if err := validateAndProcessTask(&task); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateTask(task); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONResponse(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		} else {
			writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to update task"})
		}
		return
	}

	// В случае успеха возвращаем пустой JSON-объект
	writeJSONResponse(w, http.StatusOK, map[string]string{})
}
