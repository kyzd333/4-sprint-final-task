package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

var (
	ErrWrongFormat            = errors.New("wrong data format")
	ErrIncorrectNumberOfSteps = errors.New("incorrect number of steps")
	ErrZeroDuration           = errors.New("zero duration")
)

func parsePackage(data string) (int, time.Duration, error) {
	parsed := strings.Split(data, ",")

	if len(parsed) != 2 {
		return 0, 0, ErrWrongFormat
	}

	steps, err := strconv.Atoi(parsed[0])

	if err != nil {
		return 0, 0, err
	}

	if steps <= 0 {
		return 0, 0, ErrIncorrectNumberOfSteps
	}

	duration, err := time.ParseDuration(parsed[1])

	if err != nil {
		return 0, 0, err
	}

	if duration <= 0 {
		return 0, 0, ErrZeroDuration
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)

	if err != nil {
		log.Println(err)
		return ""
	}

	distance := stepLength * float64(steps) / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	if err != nil {
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.\n",
		steps, distance, calories,
	)
}
