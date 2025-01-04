package services

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	nextDate time.Time
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

// NextDate вычисляет следующую дату для повторного выполнения задачи следуя правилам,
// которые реализуют внутренние функции.
func NextDate(now time.Time, date string, repeat string) (string, error) {
	dateTime, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	// Это присвоения необходимы для parseWeekday и parseMonth
	if now.Before(dateTime) {
		nextDate = dateTime
	} else {
		nextDate = now
	}

	elems := strings.Fields(repeat)
	if len(elems) == 0 {
		return "", errors.New("repeat is empty row")
	}

	var nextDate string

	switch elems[0] {
	case "d":
		nextDate, err = parseDays(now, dateTime, elems)
	case "y":
		nextDate, err = parseYears(now, dateTime, elems)
	case "w":
		nextDate, err = parseWeekday(elems)
	case "m":
		nextDate, err = parseMonth(elems)
	default:
		return "", errInvalidFormat
	}

	if err != nil {
		return "", err
	}

	return nextDate, nil
}

// parseDays разбирает слайс elem, чтобы получить число дней для инкриментирования даты.
func parseDays(now time.Time, date time.Time, elems []string) (string, error) {
	if len(elems) != 2 {
		return "", errInvalidFormat
	}

	days, err := strconv.Atoi(elems[1])
	if err != nil || days > 400 {
		return "", errInvalidFormat
	}

	if now.Before(date) {
		nextDate = nextDate.AddDate(0, 0, days)
	} else {
		nextDate = date
		for now.After(nextDate) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
	}

	return nextDate.Format("20060102"), nil
}

// parseYears разбирает слайс elem, чтобы получить число лет для инкриментирования даты.
func parseYears(now time.Time, date time.Time, elems []string) (string, error) {
	if len(elems) != 1 {
		return "", errInvalidFormat
	}

	if now.Before(date) {
		nextDate = nextDate.AddDate(1, 0, 0)
	} else {
		nextDate = date
		for now.After(nextDate) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
	}
	return nextDate.Format("20060102"), nil
}

// parseWeekday разбирает слайс elem, чтобы получить числа дней недели для ближашего переноса даты.
func parseWeekday(elems []string) (string, error) {
	if len(elems) != 2 {
		return "", errInvalidFormat
	}

	// Получаем список, еще непроверенных, дней
	var uncheckedDaysList []string
	if len(elems) > 1 {
		uncheckedDaysList = strings.Split(elems[1], ",")
	} else {
		return "", errInvalidFormat
	}

	weekDayList := []int{}
	// Проверка на корректность введенных дней недели
	for _, day := range uncheckedDaysList {
		num, err := strconv.Atoi(day)
		if err != nil {
			return "", errInvalidFormat
		}

		if num < 1 || num > 7 {
			return "", errInvalidFormat
		}

		weekDayList = append(weekDayList, num)
	}

	var weekName, nextWeekName string

	weekName = nextDate.Weekday().String()

	// Отсорируем список с днями недель по возрастанию (пригодится на следующем шаге)
	slices.Sort(weekDayList)

	for _, weekDay := range weekDayList {
		// Находим ближайший меньший weekDay
		if weekDay < weekStore[weekName] {
			nextWeekName = nextDate.Weekday().String()

			// Инкриментируем дату до выбранного ближайшего дня недели
			for weekStore[nextWeekName] != weekDay {
				nextDate = nextDate.AddDate(0, 0, 1)
				nextWeekName = nextDate.Weekday().String()
			}

			return nextDate.Format("20060102"), nil
		}
	}

	// Будет выполнятся в случае если в списке с днями недель все числа будут меньше базового дня.
	// Инкриментируем дату до первого в списке дня недели,
	// т.к. список сформирован в порядке возрастания (порядок проверялся этапами выше)
	for weekStore[nextWeekName] != weekDayList[0] {
		nextDate = nextDate.AddDate(0, 0, 1)
		nextWeekName = nextDate.Weekday().String()
	}

	return nextDate.Format("20060102"), nil
}

// parseMonth разбирает слайс elem, чтобы получить числа дней месяца для ближашего переноса даты.
func parseMonth(elems []string) (string, error) {
	// Получаем список, еще непроверенных, дней
	var uncheckedDaysList []string
	if len(elems) > 1 {
		uncheckedDaysList = strings.Split(elems[1], ",")
	} else {
		return "", errInvalidFormat
	}

	daysList := make([]int, len(uncheckedDaysList))

	// Создадим словарь,
	// где ключ - это номер дня (пригодится для нахождения ближайшего следующего дня месяца)
	dictDays := map[int]int{}

	// дата для нахождения последнего и предпоследнего дня текущего месяца
	tmpDate := time.Date(nextDate.Year(), nextDate.Month(), 1, 0, 0, 0, 0, time.Local)
	tmpDate = tmpDate.AddDate(0, 1, 0)

	// Проверка введенных дней
	for i, day := range uncheckedDaysList {
		dayNum, err := strconv.Atoi(day)
		if err != nil {
			return "", errInvalidFormat
		}

		if dayNum > 31 || dayNum == 0 || dayNum < -2 {
			return "", errInvalidFormat
		}

		daysList[i] = dayNum
	}

	if len(elems) == 2 {
		// Заполняем словарь с днями, учитывая -1 и -2
		for _, day := range daysList {
			switch day {
			case -1:
				lastDay := tmpDate.AddDate(0, 0, -1).Day()
				dictDays[lastDay] = 1
			case -2:
				preLastDay := tmpDate.AddDate(0, 0, -2).Day()
				dictDays[preLastDay] = 1
			default:
				dictDays[day] = 1
			}
		}

		// Ищем ближашую дату, где день совпадает с тем, что есть в словаре
		for {
			nextDate = nextDate.AddDate(0, 0, 1)
			if _, has := dictDays[nextDate.Day()]; has {
				break
			}
		}
	} else if len(elems) == 3 {
		uncheckedMonthList := strings.Split(elems[2], ",")
		monthList := make([]int, len(uncheckedMonthList))

		// Проверка введенных месяцев
		for i, month := range uncheckedMonthList {
			monthNum, err := strconv.Atoi(month)
			if err != nil {
				return "", errInvalidFormat
			}

			if monthNum < 1 || monthNum > 12 {
				return "", errInvalidFormat
			}

			monthList[i] = monthNum
		}

		// Переменная, необхожимая для нахождения минимального интервала времени между датами
		var minDuration time.Duration
		// Флаг для обозначения опроного минимума
		once := true

		// Нахождение минимального итервала
		for _, day := range daysList {
			for _, month := range monthList {
				assistDate := getAssistDate(tmpDate, day, month)

				if nextDate.Before(assistDate) {
					// Необходимо один раз обозначить опорный минимум
					if once {
						minDuration = assistDate.Sub(nextDate)
						once = false
					}

					minDuration = min(minDuration, assistDate.Sub(nextDate))
				}
			}
		}

		nextDate = nextDate.Add(minDuration).AddDate(0, 0, 1)
	} else {
		return "", errInvalidFormat
	}

	return nextDate.Format("20060102"), nil
}

// getAssistDate получяет опорную дату, необхоодимую для нахождения минимального интервала времени между датами.
// assistDate - вспомогательная дата;
// tmpDate - дата необходимая для вычисления последнего и предпоследнего дня месяца.
func getAssistDate(tmpDate time.Time, day, month int) time.Time {
	// вспомогательная переменная для нахождения минимального интервала времени между датами
	var assistDate time.Time

	switch day {
	case -1:
		// Если месяц меньше чем в actualDate, то образованная дата переносится на следующий год
		if month < int(nextDate.Month()) {
			assistDate = tmpDate.AddDate(1, 0, -1)
		} else {
			assistDate = tmpDate.AddDate(0, 0, -1)
		}

	case -2:
		if month < int(nextDate.Month()) {
			assistDate = tmpDate.AddDate(1, 0, -2)
		} else {
			assistDate = tmpDate.AddDate(0, 0, -2)
		}
	default:
		if month < int(nextDate.Month()) {
			assistDate = time.Date(nextDate.Year()+1, time.Month(month), day, 0, 0, 0, 0, time.Local)
		} else {
			assistDate = time.Date(nextDate.Year(), time.Month(month), day, 0, 0, 0, 0, time.Local)
		}
	}

	return assistDate
}
