package db

import (
    "database/sql"
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

// GetTask получает одну задачу из БД по её ID.
func GetTask(id string) (Task, error) {
	var task Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	// QueryRow идеально подходит для получения одной записи.
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		// Ошибка sql.ErrNoRows означает, что задача не найдена.
		// Мы вернем эту ошибку, а API-слой обработает её и выдаст 404.
		return Task{}, err
	}

	return task, nil
}

// UpdateTask обновляет существующую задачу в БД.
func UpdateTask(task Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	// Проверяем, была ли действительно обновлена строка.
	// Если нет, значит, задачи с таким ID не существует.
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // Используем ту же ошибку, что и при поиске, для консистентности.
	}

	return nil
}

// DeleteTask удаляет задачу из БД по её ID.
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	// Проверяем, была ли действительно удалена строка.
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // Задача с таким ID не найдена.
	}

	return nil
}

// UpdateTaskDate обновляет только дату у существующей задачи.
func UpdateTaskDate(id string, newDate string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	res, err := db.Exec(query, newDate, id)
	if err != nil {
		return err
	}

	// Проверяем, была ли действительно обновлена строка.
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // Задача с таким ID не найдена.
	}

	return nil
}
