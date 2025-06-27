package api

import (
	"encoding/json"
	"errors"
	"go1f/pkg/db"
	"net/http"
	"time"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	// 1. Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	// 2. Валидация данных
	if err := validateAndProcessTask(&task); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// 3. Добавляем задачу в базу данных
	id, err := db.AddTask(task)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to add task"})
		return
	}

	// 4. Отправляем успешный ответ с ID новой задачи
	writeJSONResponse(w, http.StatusCreated, map[string]int64{"id": id})
}

// validateAndProcessTask проверяет корректность данных задачи и обрабатывает дату.
func validateAndProcessTask(task *db.Task) error {
	// Проверка на обязательное поле title
	if task.Title == "" {
		return errors.New("title is required")
	}

	now := time.Now()
	// Если дата не указана, используем сегодняшнюю
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	// Проверяем формат даты
	taskDate, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return errors.New("invalid date format, expected YYYYMMDD")
	}

	// Если указано правило повторения, оно должно быть валидным
	if task.Repeat != "" {
		// Используем NextDate для проверки правила. Нам не важен результат, только ошибка.
		_, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err // Возвращаем ошибку от NextDate (e.g., "unsupported rule")
		}
	}

	// Если дата в прошлом
	if taskDate.Before(now.Truncate(24 * time.Hour)) {
		if task.Repeat == "" {
			// без повторения - ставим на сегодня
			task.Date = now.Format(DateFormat)
		} else {
			// с повторением - вычисляем следующую дату
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err // Хотя мы уже проверяли, на всякий случай
			}
			task.Date = next
		}
	}

	return nil
}
