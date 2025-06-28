package api

import "net/http"

// taskHandler является диспетчером для всех запросов к /api/task.
// Он определяет HTTP-метод и вызывает соответствующий обработчик.
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	// DELETE будет добавлен на следующем шаге
	default:
		writeJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}
