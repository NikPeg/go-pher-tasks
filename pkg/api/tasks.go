package api

import (
	"go1f/pkg/db"
	"net/http"
	"log"
)

// TasksResponse определяет структуру JSON-ответа для списка задач.
type TasksResponse struct {
	Tasks []db.Task `json:"tasks"`
}

// tasksHandler обрабатывает GET-запросы к /api/tasks.
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	// Получаем необязательный параметр search из URL
	searchQuery := r.URL.Query().Get("search")

	const tasksLimit = 50

	// Передаем поисковый запрос в функцию GetTasks
	tasks, err := db.GetTasks(searchQuery, tasksLimit)
	if err != nil {
		// Для отладки можно выводить ошибку в лог
		log.Printf("Error getting tasks: %v", err)
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to get tasks"})
		return
	}

	response := TasksResponse{Tasks: tasks}
	writeJSONResponse(w, http.StatusOK, response)
}
