package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// DateFormat определяет стандартный формат даты для всего проекта.
const DateFormat = "20060102"

// NextDate вычисляет следующую дату на основе правила повторения.
func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat rule cannot be empty")
	}

	// Парсим исходную дату задачи
	startDate, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %w", err)
	}

	parts := strings.Split(repeat, " ")
	ruleType := parts[0]

	currentDate := startDate

	switch ruleType {
	case "y":
		// Ежегодное повторение
		for {
			currentDate = currentDate.AddDate(1, 0, 0)
			if currentDate.After(now) {
				return currentDate.Format(DateFormat), nil
			}
		}

	case "d":
		// Повторение через N дней
		if len(parts) < 2 {
			return "", errors.New("missing number of days for 'd' rule")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("invalid number of days: %s", parts[1])
		}
		if days <= 0 || days > 400 {
			return "", fmt.Errorf("number of days must be between 1 and 400, got %d", days)
		}

		for {
			currentDate = currentDate.AddDate(0, 0, days)
			if currentDate.After(now) {
				return currentDate.Format(DateFormat), nil
			}
		}

	default:
		// Для начала считаем все остальные правила неподдерживаемыми
		return "", fmt.Errorf("unsupported repeat rule: %s", ruleType)
	}
}

// nextDateHandler обрабатывает GET-запросы к /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из URL
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		// Если 'now' не указан, берем текущее время
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, `invalid 'now' date format, expected YYYYMMDD`, http.StatusBadRequest)
			return
		}
	}

	if dateStr == "" || repeat == "" {
		http.Error(w, `'date' and 'repeat' parameters are required`, http.StatusBadRequest)
		return
	}

	// Вызываем нашу основную функцию
	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}
