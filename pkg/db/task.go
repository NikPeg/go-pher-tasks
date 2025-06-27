package db

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
// limit - максимальное количество возвращаемых задач.
func GetTasks(limit int) ([]Task, error) {
	// SQL-запрос для выбора задач с сортировкой по дате и ограничением количества
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`

	// Выполняем запрос
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Важно закрыть rows после использования

	// Инициализируем пустой, но не nil-слайс.
	// Это гарантирует, что при отсутствии задач мы получим `[]`, а не `null` в JSON.
	tasks := make([]Task, 0)

	// Итерируемся по результатам запроса
	for rows.Next() {
		var t Task
		// Сканируем каждую строку в структуру Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	// Проверяем на ошибки, которые могли возникнуть во время итерации
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
