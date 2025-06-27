// main.go

package main

import (
	"go1f/pkg/server" // Импортируем наш пакет server, используя алиас из go.mod
	"log"
)

func main() {
	// Запускаем сервер и, если возникнет ошибка, логируем ее и завершаем программу
	if err := server.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
