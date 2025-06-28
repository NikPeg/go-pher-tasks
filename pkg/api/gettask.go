package api

import (
	"database/sql"
	"errors"
	"go1f/pkg/db"
	"net/http"
	"strconv"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "id parameter is required"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONResponse(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		} else {
			writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to get task"})
		}
		return
	}

	// Конвертируем db.Task в APITask для корректного JSON-ответа (с id: "string")
	apiTask := APITask{
		ID:      strconv.FormatInt(task.ID, 10),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	writeJSONResponse(w, http.StatusOK, apiTask)
}
