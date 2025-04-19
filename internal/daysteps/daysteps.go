package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	sc "github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

// Основные виды ошибок.
var (
	ErrIncorrectNumberOfArgs = errors.New("неверное количество параметров")
	ErrStepsNumber           = errors.New("количество шагов должно быть больше нуля")
	ErrDuration              = errors.New("неверная длительность")
)

// Основыне виды констант.
const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// Переводит данные о тренировке в переменные соответствующих типов.
func parsePackage(data string) (int, time.Duration, error) {
	if data == "" {
		return 0, 0, ErrIncorrectNumberOfArgs
	}
	components := strings.Split(data, ",")
	if len(components) != 2 {
		return 0, 0, ErrIncorrectNumberOfArgs
	}
	steps, err := strconv.Atoi(components[0])
	if err != nil {
		return 0, 0, err
	}
	if steps <= 0 {
		return 0, 0, ErrStepsNumber
	}
	duration, err := time.ParseDuration(components[1])
	if duration <= 0 {
		return 0, 0, ErrDuration
	}
	if err != nil {
		return 0, 0, err
	}
	return steps, duration, nil
}

// DayActionInfo выводит данные об активности за день.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	if steps <= 0 {
		log.Println(ErrStepsNumber.Error())
		return ""
	}
	distance := float64(steps) * stepLength
	distance = distance / mInKm
	calories, err := sc.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps,
		distance,
		calories)
}
