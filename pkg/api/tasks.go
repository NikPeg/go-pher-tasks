package api

import (
	"go1f/pkg/db"
	"net/http"
	"strconv"
)

// APITask - это структура для представления задачи в JSON-ответе API.
// Поле ID здесь является строкой, как того требуют тесты.
type APITask struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// TasksResponse определяет структуру JSON-ответа для списка задач.
// Теперь она использует APITask вместо db.Task.
type TasksResponse struct {
	Tasks []APITask `json:"tasks"`
}

// tasksHandler обрабатывает GET-запросы к /api/tasks.
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	searchQuery := r.URL.Query().Get("search")
	const tasksLimit = 50

	// 1. Получаем задачи из БД в "родном" формате (с ID типа int64)
	dbTasks, err := db.GetTasks(searchQuery, tasksLimit)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to get tasks"})
		return
	}

	// 2. Конвертируем задачи из формата БД в формат API
	apiTasks := make([]APITask, 0, len(dbTasks))
	for _, task := range dbTasks {
		apiTasks = append(apiTasks, APITask{
			ID:      strconv.FormatInt(task.ID, 10),
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})
	}

	// 3. Формируем и отправляем успешный ответ
	response := TasksResponse{Tasks: apiTasks}
	writeJSONResponse(w, http.StatusOK, response)
}
