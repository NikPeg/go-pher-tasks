package api

import (
	"database/sql"
	"errors"
	"go1f/pkg/db"
	"net/http"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "id parameter is required"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONResponse(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		} else {
			writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete task"})
		}
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{})
}
