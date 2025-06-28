package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite" // Регистрируем драйвер SQLite
)

// db является глобальной переменной для хранения подключения к базе данных.
var db *sql.DB

// schema определяет структуру базы данных (таблица и индекс).
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL,
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128) NOT NULL
);
CREATE INDEX IF NOT EXISTS date_idx ON scheduler (date);
`

// Init инициализирует подключение к базе данных.
// Если файл БД не существует, он будет создан вместе с необходимой таблицей и индексом.
func Init(dbFile string) error {
	// Проверяем, существует ли файл базы данных.
	_, err := os.Stat(dbFile)
	// os.IsNotExist(err) вернет true, если файл не найден.
	install := os.IsNotExist(err)

	// Открываем (или создаем) файл базы данных.
	conn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	db = conn // Сохраняем подключение в глобальную переменную

	// Если база данных только что была создана (install == true),
	// выполняем SQL-запрос для создания таблиц.
	if install {
		if _, err := db.Exec(schema); err != nil {
			return err
		}
	}

	return nil
}

// AddTask добавляет новую задачу в базу данных.
// Возвращает ID созданной задачи.
func AddTask(task Task) (int64, error) {
	// SQL-запрос для вставки данных. Используем '?' как плейсхолдеры для безопасности.
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	// Выполняем запрос
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	// Получаем ID последней вставленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
