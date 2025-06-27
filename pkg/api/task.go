package api

import "net/http"

// taskHandler является диспетчером для всех запросов к /api/task.
// Он определяет HTTP-метод и вызывает соответствующий обработчик.
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	// Другие методы (GET, PUT, DELETE) будут добавлены на следующих шагах
	default:
		// Если метод не поддерживается, возвращаем ошибку
		writeJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}
