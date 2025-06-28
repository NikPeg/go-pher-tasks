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

// --- Основная публичная функция ---

// NextDate вычисляет следующую дату на основе правила повторения.
// Она выступает в роли диспетчера, вызывая нужный обработчик правила.
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

	var nextDate time.Time

	switch ruleType {
	case "y":
		nextDate, err = nextYearlyDate(startDate, now)
	case "d":
		nextDate, err = nextDailyDate(startDate, now, parts)
	case "w":
		nextDate, err = nextWeeklyDate(startDate, now, parts)
	case "m":
		nextDate, err = nextMonthlyDate(startDate, now, parts)
	default:
		return "", fmt.Errorf("unsupported repeat rule: %s", ruleType)
	}

	if err != nil {
		return "", err
	}

	return nextDate.Format(DateFormat), nil
}

// --- Функции-обработчики для каждого правила ---

// nextYearlyDate ищет следующую годовую дату.
func nextYearlyDate(current, now time.Time) (time.Time, error) {
	for {
		current = current.AddDate(1, 0, 0)
		if current.After(now) {
			return current, nil
		}
	}
}

// nextDailyDate ищет следующую дату с шагом в N дней.
func nextDailyDate(current, now time.Time, parts []string) (time.Time, error) {
	if len(parts) < 2 {
		return time.Time{}, errors.New("missing number of days for 'd' rule")
	}
	days, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid number of days: %s", parts[1])
	}
	if days <= 0 || days > 400 {
		return time.Time{}, fmt.Errorf("number of days must be between 1 and 400, got %d", days)
	}

	for {
		current = current.AddDate(0, 0, days)
		if current.After(now) {
			return current, nil
		}
	}
}

// nextWeeklyDate ищет следующую дату, подходящую под правило дней недели.
func nextWeeklyDate(current, now time.Time, parts []string) (time.Time, error) {
	allowedDays, err := parseWeeklyRule(parts)
	if err != nil {
		return time.Time{}, err
	}

	for i := 0; i < 365*2; i++ { // Ограничим поиск двумя годами
		current = current.AddDate(0, 0, 1)
		if allowedDays[current.Weekday()] && current.After(now) {
			return current, nil
		}
	}
	return time.Time{}, errors.New("could not find a matching date for 'w' rule within 2 years")
}

// nextMonthlyDate ищет следующую дату, подходящую под правило дней и месяцев.
func nextMonthlyDate(current, now time.Time, parts []string) (time.Time, error) {
	rule, err := parseMonthlyRule(parts)
	if err != nil {
		return time.Time{}, err
	}

	for i := 0; i < 365*10; i++ { // Ограничим поиск 10 годами
		current = current.AddDate(0, 0, 1)

		// Пропускаем, если месяц не подходит
		if !rule.allMonthsAllowed && !rule.allowedMonths[current.Month()] {
			continue
		}

		// Проверяем, подходит ли день
		if isDayMatch(current, rule.allowedDays) && current.After(now) {
			return current, nil
		}
	}
	return time.Time{}, errors.New("could not find a matching date for 'm' rule within 10 years")
}

// --- Вспомогательные функции-парсеры и проверки ---

// parseWeeklyRule разбирает строковое правило для дней недели.
func parseWeeklyRule(parts []string) (map[time.Weekday]bool, error) {
	if len(parts) < 2 {
		return nil, errors.New("missing days of week for 'w' rule")
	}
	allowedDays := make(map[time.Weekday]bool)
	weekdaysStr := strings.Split(parts[1], ",")
	for _, dayStr := range weekdaysStr {
		day, err := strconv.Atoi(dayStr)
		if err != nil || day < 1 || day > 7 {
			return nil, fmt.Errorf("invalid day of week: %s", dayStr)
		}
		allowedDays[time.Weekday(day%7)] = true // 1..6 -> Mon..Sat, 7 -> Sun (0)
	}
	return allowedDays, nil
}

// monthlyRule хранит разобранное правило для месяцев.
type monthlyRule struct {
	allowedDays      map[int]bool
	allowedMonths    map[time.Month]bool
	allMonthsAllowed bool
}

// parseMonthlyRule разбирает строковое правило для дней и месяцев.
func parseMonthlyRule(parts []string) (monthlyRule, error) {
	var rule monthlyRule
	if len(parts) < 2 {
		return rule, errors.New("missing day numbers for 'm' rule")
	}

	// Разбираем дни
	rule.allowedDays = make(map[int]bool)
	daysStr := strings.Split(parts[1], ",")
	for _, dayStr := range daysStr {
		day, err := strconv.Atoi(dayStr)
		if err != nil || ((day < 1 || day > 31) && (day != -1 && day != -2)) {
			return rule, fmt.Errorf("invalid day number: %s", dayStr)
		}
		rule.allowedDays[day] = true
	}

	// Разбираем месяцы
	rule.allMonthsAllowed = true
	if len(parts) > 2 {
		rule.allMonthsAllowed = false
		rule.allowedMonths = make(map[time.Month]bool)
		monthsStr := strings.Split(parts[2], ",")
		for _, monthStr := range monthsStr {
			month, err := strconv.Atoi(monthStr)
			if err != nil || month < 1 || month > 12 {
				return rule, fmt.Errorf("invalid month number: %s", monthStr)
			}
			rule.allowedMonths[time.Month(month)] = true
		}
	}
	return rule, nil
}

// isDayMatch проверяет, соответствует ли дата правилам дня для правила "m".
func isDayMatch(date time.Time, allowedDays map[int]bool) bool {
	dayOfMonth := date.Day()
	if allowedDays[dayOfMonth] {
		return true
	}
	if allowedDays[-1] {
		lastDayOfMonth := date.AddDate(0, 1, -date.Day()).Day()
		if dayOfMonth == lastDayOfMonth {
			return true
		}
	}
	if allowedDays[-2] {
		secondToLastDayOfMonth := date.AddDate(0, 1, -date.Day()).Day() - 1
		if dayOfMonth == secondToLastDayOfMonth {
			return true
		}
	}
	return false
}

// --- Обработчик HTTP-запроса (без изменений) ---

// nextDateHandler обрабатывает GET-запросы к /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
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

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}
