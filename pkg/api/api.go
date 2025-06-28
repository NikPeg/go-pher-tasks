package api

import "net/http"

func Init() {
	// Незащищенные эндпоинты
	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/nextdate", nextDateHandler)

	// Защищенные эндпоинты
	http.Handle("/api/task", authMiddleware(http.HandlerFunc(taskHandler)))
	http.Handle("/api/tasks", authMiddleware(http.HandlerFunc(tasksHandler)))
	http.Handle("/api/task/done", authMiddleware(http.HandlerFunc(doneTaskHandler)))
}
