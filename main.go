package main

import (
	"go1f/pkg/db"
	"go1f/pkg/server"
	"log"
	"os"
)

const defaultDBFile = "scheduler.db"

func main() {
	// Получаем путь к файлу БД из переменной окружения TODO_DBFILE.
	// Если она не установлена, используем значение по умолчанию.
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDBFile
	}

	// Инициализируем базу данных ПЕРЕД запуском сервера.
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// Запускаем веб-сервер.
	if err := server.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
