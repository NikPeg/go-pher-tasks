package server

import (
	"go1f/pkg/api"
	"log"
	"net/http"
	"os"
)

const defaultPort = "7540"
const webDir = "./web" // Директория с файлами фронтенда

// Run запускает веб-сервер
func Run() error {
	// Получаем порт из переменной окружения TODO_PORT, иначе используем порт по умолчанию
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	api.Init()

	// Регистрируем файловый сервер для раздачи статических файлов
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Server starting on http://localhost:%s", port)

	return http.ListenAndServe(":"+port, nil)
}
