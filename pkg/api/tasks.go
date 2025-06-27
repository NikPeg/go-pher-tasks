package api

import (
	"go1f/pkg/db"
	"net/http"
)

// TasksResponse определяет структуру JSON-ответа для списка задач.
type TasksResponse struct {
	Tasks []db.Task `json:"tasks"`
}

// tasksHandler обрабатывает GET-запросы к /api/tasks.
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// Мы обрабатываем только GET-запросы на этом эндпоинте
	if r.Method != http.MethodGet {
		writeJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	const tasksLimit = 50 // Ограничим выборку 50 задачами

	// Получаем задачи из базы данных
	tasks, err := db.GetTasks(tasksLimit)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to get tasks"})
		return
	}

	// Формируем и отправляем успешный ответ
	response := TasksResponse{Tasks: tasks}
	writeJSONResponse(w, http.StatusOK, response)
}
