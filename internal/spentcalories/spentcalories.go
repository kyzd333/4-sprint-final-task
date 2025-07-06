package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

var (
	ErrWrongFormat            = errors.New("wrong data format")
	ErrIncorrectNumberOfSteps = errors.New("incorrect number of steps")
	ErrZeroDuration           = errors.New("zero duration")
	ErrIncorrectBodyInfo      = errors.New("incorrect body info")
	ErrUnknownTrainingType    = errors.New("неизвестный тип тренировки")
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parsed := strings.Split(data, ",")

	if len(parsed) != 3 {
		return 0, "", 0, ErrWrongFormat
	}

	steps, err := strconv.Atoi(parsed[0])

	if err != nil {
		return 0, "", 0, err
	}

	if steps <= 0 {
		return 0, "", 0, ErrIncorrectNumberOfSteps
	}

	duration, err := time.ParseDuration(parsed[2])

	if err != nil {
		return 0, "", 0, err
	}

	if duration <= 0 {
		return 0, "", 0, ErrZeroDuration
	}

	return steps, parsed[1], duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient

	return float64(steps) * stepLength / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, typeOfTraining, duration, err := parseTraining(data)

	if err != nil {
		log.Println(err)
	}

	if weight <= 0 || height <= 0 {
		return "", ErrIncorrectBodyInfo
	}

	var calories float64
	switch typeOfTraining {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", ErrUnknownTrainingType
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		typeOfTraining, duration.Hours(), distance(steps, height),
		meanSpeed(steps, height, duration), calories,
	), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, ErrIncorrectNumberOfSteps
	}

	if duration <= 0 {
		return 0, ErrZeroDuration
	}

	if weight <= 0 || height <= 0 {
		return 0, ErrIncorrectBodyInfo
	}

	speed := meanSpeed(steps, height, duration)

	return (weight * speed * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := RunningSpentCalories(steps, weight, height, duration)

	if err != nil {
		return 0, err
	}

	return calories * walkingCaloriesCoefficient, nil
}
