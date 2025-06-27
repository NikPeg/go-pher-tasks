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

	startDate, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %w", err)
	}

	parts := strings.Split(repeat, " ")
	ruleType := parts[0]

	currentDate := startDate

	switch ruleType {
	case "y":
		// Увеличиваем дату на год, пока она не станет больше 'now'
		for {
			currentDate = currentDate.AddDate(1, 0, 0)
			if currentDate.After(now) {
				// Учтем случай с високосным годом, например, 29 февраля
				// Если исходная дата была 29.02, а следующий год не високосный, time.AddDate вернет 28.02
                // Чтобы исправить это и всегда получать 1 марта, можно добавить проверку, но для простоты оставим поведение time.AddDate.
				return currentDate.Format(DateFormat), nil
			}
		}

	case "d":
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
		// Увеличиваем дату на N дней, пока она не станет больше 'now'
		for {
			currentDate = currentDate.AddDate(0, 0, days)
			if currentDate.After(now) {
				return currentDate.Format(DateFormat), nil
			}
		}

	case "w":
		// === Шаг 1: Разбор правила ===
		if len(parts) < 2 {
			return "", errors.New("missing days of week for 'w' rule")
		}
		weekdaysStr := strings.Split(parts[1], ",")
		allowedDays := make(map[time.Weekday]bool)
		for _, dayStr := range weekdaysStr {
			day, err := strconv.Atoi(dayStr)
			if err != nil || day < 1 || day > 7 {
				return "", fmt.Errorf("invalid day of week: %s", dayStr)
			}
			// 1(пн)->1, ..., 6(сб)->6, 7(вс)->0
			weekday := time.Weekday(day % 7)
			allowedDays[weekday] = true
		}

		// === Шаг 2: Поиск подходящей даты ===
		for i := 0; i < 365*2; i++ { // Ограничим поиск двумя годами
			currentDate = currentDate.AddDate(0, 0, 1)
			if allowedDays[currentDate.Weekday()] {
				if currentDate.After(now) {
					return currentDate.Format(DateFormat), nil
				}
			}
		}
		return "", errors.New("could not find a matching date for 'w' rule within 2 years")


	case "m":
		// === Шаг 1: Разбор правила ===
		if len(parts) < 2 {
			return "", errors.New("missing day numbers for 'm' rule")
		}

		allowedDays := make(map[int]bool)
		daysStr := strings.Split(parts[1], ",")
		for _, dayStr := range daysStr {
			day, err := strconv.Atoi(dayStr)
			if err != nil {
				return "", fmt.Errorf("invalid day number: %s", dayStr)
			}
			if (day < 1 || day > 31) && (day != -1 && day != -2) {
				return "", fmt.Errorf("day number must be between 1-31, or -1, -2. Got: %d", day)
			}
			allowedDays[day] = true
		}

		allMonthsAllowed := true
		allowedMonths := make(map[time.Month]bool)
		if len(parts) > 2 {
			allMonthsAllowed = false
			monthsStr := strings.Split(parts[2], ",")
			for _, monthStr := range monthsStr {
				month, err := strconv.Atoi(monthStr)
				if err != nil {
					return "", fmt.Errorf("invalid month number: %s", monthStr)
				}
				if month < 1 || month > 12 {
					return "", fmt.Errorf("month number must be between 1 and 12. Got: %d", month)
				}
				allowedMonths[time.Month(month)] = true
			}
		}

		// === Шаг 2: Поиск подходящей даты ===
		for i := 0; i < 365*10; i++ { // Ограничим цикл 10 годами
			currentDate = currentDate.AddDate(0, 0, 1)

			if !allMonthsAllowed && !allowedMonths[currentDate.Month()] {
				continue
			}

			dateIsValid := false
			dayOfMonth := currentDate.Day()

			if allowedDays[dayOfMonth] {
				dateIsValid = true
			}
			if allowedDays[-1] {
				lastDayOfMonth := currentDate.AddDate(0, 1, -currentDate.Day()).Day()
				if dayOfMonth == lastDayOfMonth {
					dateIsValid = true
				}
			}
			if allowedDays[-2] {
				secondToLastDayOfMonth := currentDate.AddDate(0, 1, -currentDate.Day()).Day() - 1
				if dayOfMonth == secondToLastDayOfMonth {
					dateIsValid = true
				}
			}

			if dateIsValid && currentDate.After(now) {
				return currentDate.Format(DateFormat), nil
			}
		}
		return "", errors.New("could not find a matching date for 'm' rule within 10 years")

	default:
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
