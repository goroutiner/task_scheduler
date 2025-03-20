package services

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	// Словарь связывающий навзания недель с их порядковым номером
	weekStore = map[string]int{
		"Monday":    1,
		"Tuesday":   2,
		"Wednesday": 3,
		"Thursday":  4,
		"Friday":    5,
		"Saturday":  6,
		"Sunday":    7,
	}

	errInvalidFormat = errors.New("invalid `repeat` format")
)

// GetNextDate вычисляет следующую дату относительно заданной, в соответствии в правилом повторения.
func (s *TaskService) GetNextDate(now time.Time, date string, repeat string) (string, error) {
	var nextDate string

	dateTime, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	elems := strings.Fields(repeat)
	if len(elems) == 0 {
		return "", errInvalidFormat
	}

	switch elems[0] {
	case "d":
		nextDate, err = nextDateByDay(now, dateTime, elems)
	case "y":
		nextDate, err = nextDateByYear(now, dateTime)
	case "w":
		nextDate, err = nextDateByWeekday(now, dateTime, elems)
	case "m":
		nextDate, err = nextDateByDayOfMonth(now, dateTime, elems)
	default:
		return "", errInvalidFormat
	}

	return nextDate, err
}

// nextDateByDay получает следующую дату, инкриментируя по дням.
func nextDateByDay(now time.Time, date time.Time, elems []string) (string, error) {
	if len(elems) != 2 {
		return "", errInvalidFormat
	}

	dayInc, err := strconv.Atoi(elems[1])
	if err != nil || dayInc < 1 || dayInc > 366 {
		return "", errInvalidFormat
	}

	if date.Before(now) {
		for date.Before(now) {
			date = date.AddDate(0, 0, dayInc)
		}
	} else {
		date = date.AddDate(0, 0, dayInc)
	}

	return date.Format("20060102"), nil
}

// nextDateByYear получает следующую дату, инкриментируя по годам.
func nextDateByYear(now time.Time, date time.Time) (string, error) {
	if date.Before(now) {
		for date.Before(now) {
			date = date.AddDate(1, 0, 0)
		}
	} else {
		date = date.AddDate(1, 0, 0)
	}

	return date.Format("20060102"), nil
}

// nextDateByWeekday получает следующую дату, в соответствии с днем недели, инкриментируя по дням.
func nextDateByWeekday(now time.Time, date time.Time, elems []string) (string, error) {
	if len(elems) != 2 {
		return "", errInvalidFormat
	}

	if date.Before(now) {
		date = now
	}

	// Получаем список порядковых номеров дней недели
	daysOfWeakList := strings.Split(elems[1], ",")
	dowDir := make(map[int]bool, len(daysOfWeakList))
	// Проверка на корректность введенных дней недели
	for _, dow := range daysOfWeakList {
		num, err := strconv.Atoi(dow)
		if err != nil || num < 1 || num > 7 {
			return "", errInvalidFormat
		}

		dowDir[num] = true
	}

	for {
		date = date.AddDate(0, 0, 1)
		nameOfDay := date.Weekday().String()
		if dowDir[weekStore[nameOfDay]] {
			break
		}
	}

	return date.Format("20060102"), nil
}

// nextDateByDayOfMonth получает следующую дату, в соответствии с днем месяца, инкриментируя по дням.
func nextDateByDayOfMonth(now time.Time, date time.Time, elems []string) (string, error) {
	if len(elems) == 1 || len(elems) > 3 {
		return "", errInvalidFormat
	}

	if date.Before(now) {
		date = now
	}

	daysList := strings.Split(elems[1], ",")
	dayDir := make(map[int]bool, len(daysList))

	for _, day := range daysList {
		num, err := strconv.Atoi(day)
		if err != nil || num > 31 || num == 0 || num < -2 {
			return "", errInvalidFormat
		}

		dayDir[num] = true
	}

	var monthList []string
	monthDir := make(map[int]bool)

	if len(elems) == 3 {
		monthList = strings.Split(elems[2], ",")
		for _, month := range monthList {
			num, err := strconv.Atoi(month)
			if err != nil || num < 1 || num > 12 {
				return "", errInvalidFormat
			}

			monthDir[num] = true
		}
	}

	for {
		date = date.AddDate(0, 0, 1)
		negDay := getNegativeDay(date)
		if len(elems) == 2 && (dayDir[date.Day()] || dayDir[negDay]) {
			break
		}

		if len(elems) == 3 && (dayDir[date.Day()] || dayDir[negDay]) && monthDir[int(date.Month())] {
			break
		}
	}

	return date.Format("20060102"), nil
}

// getNegativeDay вспомогательная функция для получения отрицательного дня по положительному
// Принцип работы на примере февраля (год не високосный): числа месяца 1, 2 ...  27, 28 соотносятся попарно по порядку,
// но с инверсией 1 -> 28(-1), 2 -> 27(-2)
func getNegativeDay(date time.Time) int {
	tmpDate := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.Local)
	return date.Day() - tmpDate.Day() - 1
}
