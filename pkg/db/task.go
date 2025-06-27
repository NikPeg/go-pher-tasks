package db

import (
	"fmt"
	"time"
)

// Task представляет собой одну задачу в планировщике.
// Теги json используются для сериализации/десериализации данных при обмене с фронтендом.
type Task struct {
	ID      int64  `json:"id,omitempty"` // omitempty, чтобы не включать в JSON, если значение нулевое
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// GetTasks извлекает из БД список задач, отсортированных по дате.
// Поддерживает опциональный поиск по тексту или по дате.
// search - строка для поиска.
// limit - максимальное количество возвращаемых задач.
func GetTasks(search string, limit int) ([]Task, error) {
	baseQuery := `SELECT id, date, title, comment, repeat FROM scheduler`
	whereClause := ""
	args := []interface{}{}

	if search != "" {
		// Проверяем, является ли поисковой запрос датой в формате ДД.ММ.ГГГГ
		if t, err := time.Parse("02.01.2006", search); err == nil {
			// Это поиск по дате
			whereClause = " WHERE date = ?"
			// Преобразуем дату в формат, который хранится в БД (ГГГГММДД)
			dateForDB := t.Format("20060102")
			args = append(args, dateForDB)
		} else {
			// Это поиск по тексту
			// Используем LOWER() для регистронезависимого поиска
			whereClause = " WHERE LOWER(title) LIKE LOWER(?) OR LOWER(comment) LIKE LOWER(?)"
			// Добавляем wildcards '%' для поиска подстроки
			searchTerm := fmt.Sprintf("%%%s%%", search)
			args = append(args, searchTerm, searchTerm)
		}
	}

	// Собираем финальный запрос
	finalQuery := fmt.Sprintf("%s%s ORDER BY date ASC LIMIT ?", baseQuery, whereClause)
	args = append(args, limit)

	// Выполняем запрос
	rows, err := db.Query(finalQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
