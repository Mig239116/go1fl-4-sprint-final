package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные виды ошибок.
var (
	ErrIncorrectNumberOfArgs = errors.New("неверное количество параметров")
	ErrStepsNumber           = errors.New("количество шагов должно быть больше нуля")
	ErrHeight                = errors.New("неверный рост")
	ErrWeight                = errors.New("неверный вес")
	ErrDuration              = errors.New("неверная длительность")
	ErrActivity              = errors.New("неизвестный тип тренировки")
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// Обрабатывает строку с данными о тренировке и возвращает в виде соответствующих типов.
func parseTraining(data string) (int, string, time.Duration, error) {
	components := strings.Split(data, ",")
	if len(components) != 3 {
		return 0, "", 0, ErrIncorrectNumberOfArgs
	}
	steps, err := strconv.Atoi(components[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, ErrStepsNumber
	}
	duration, err := time.ParseDuration(components[2])
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, ErrDuration
	}
	if components[1] != "Бег" && components[1] != "Ходьба" {
		return 0, "", 0, ErrActivity
	}
	return steps, components[1], duration, nil
}

// Переводит количество шагов в дистанцию в км.
func distance(steps int, height float64) float64 {
	if steps <= 0 {
		return 0
	}
	return height * stepLengthCoefficient * float64(steps) / mInKm
}

// Рассчитывает среднюю скорость.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	return dist / duration.Hours()
}

// TrainingInfo возвращает данные о тренировке.
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}
	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	var calories float64
	var calcErr error
	switch activity {
	case "Бег":
		calories, calcErr = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, calcErr = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", ErrActivity
	}
	if calcErr != nil {
		return "", err
	}
	return fmt.Sprintf(`Тип тренировки: %s
Длительность: %.2f ч.
Дистанция: %.2f км.
Скорость: %.2f км/ч
Сожгли калорий: %.2f
`,
		activity,
		duration.Hours(),
		dist,
		speed,
		calories), nil
}

// RunnningSpentCalories рассчитывает калории при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := spentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}
	return calories, nil
}

//WalkingSpentCalories рассчитывает калории при ходьбе.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := spentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}
	return calories * walkingCaloriesCoefficient, nil
}

//Производит расчет калорий.
func spentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, ErrStepsNumber
	}
	if weight <= 0 {
		return 0, ErrWeight
	}
	if height <= 0 {
		return 0, ErrHeight
	}
	if duration <= 0 {
		return 0, ErrDuration
	}
	speed := meanSpeed(steps, height, duration)
	return (weight * speed * duration.Minutes()) / minInH, nil
}
