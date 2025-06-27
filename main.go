package main

import (
	"log"
	"net/http"
	"os"
)

// Порт по умолчанию
const defaultPort = "7540"

func main() {
	// Получаем порт из переменной окружения, или используем значение по умолчанию
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	// Указываем, что все запросы к серверу нужно обрабатывать
	// с помощью файлового сервера, который работает с файлами из директории "web"
	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Printf("Server starting on port %s...\n", port)

	// Запускаем сервер
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
