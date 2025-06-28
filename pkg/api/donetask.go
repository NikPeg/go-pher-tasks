package api

import (
	"database/sql"
	"errors"
	"go1f/pkg/db"
	"net/http"
	"time"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Этот эндпоинт принимает только POST запросы
	if r.Method != http.MethodPost {
		writeJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "id parameter is required"})
		return
	}

	// 1. Получаем задачу, чтобы узнать ее правило повторения
	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONResponse(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		} else {
			writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to get task"})
		}
		return
	}

	// 2. Определяем, удалить задачу или перенести на новую дату
	if task.Repeat == "" {
		// Одноразовая задача - удаляем
		if err := db.DeleteTask(id); err != nil {
			// Ошибка здесь маловероятна, так как мы только что получили задачу,
			// но лучше обработать для надежности.
			writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete task"})
			return
		}
	} else {
		// Повторяющаяся задача - вычисляем и устанавливаем новую дату
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "failed to calculate next date: " + err.Error()})
			return
		}
		if err := db.UpdateTaskDate(id, nextDate); err != nil {
			writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to update task date"})
			return
		}
	}

	// 3. Отправляем успешный ответ
	writeJSONResponse(w, http.StatusOK, map[string]string{})
}
