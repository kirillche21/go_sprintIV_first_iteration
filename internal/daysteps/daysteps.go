package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength float64 = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	split := strings.Split(data, ",")
	if len(split) != 2 {
		return 0, 0, fmt.Errorf("Длина слайса %d != 2", len(split))
	}

	valueOfSteps, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, 0, fmt.Errorf("Ошибка преобразования 1-го элемента: %s", err)
	}
	if valueOfSteps <= 0 {
		return 0, 0, fmt.Errorf("Ошибка: Количество шагов %d <= 0", valueOfSteps)
	}

	durationOfWalk, err := time.ParseDuration(split[1])
	if err != nil {
		return 0, 0, fmt.Errorf("Ошибка преобразования 2-го элемента: %s", err)
	}
	if durationOfWalk <= 0 {
		return 0, 0, fmt.Errorf("Ошибка: длительность прогулки < 0 %d", durationOfWalk)
	}

	return valueOfSteps, durationOfWalk, nil
}

func DayActionInfo(data string, weight, height float64) string {
	valueOfSteps, durationOfWalk, err := parsePackage(data)
	if err != nil {
		return ""
	}

	if valueOfSteps <= 0 {
		return ""
	}

	distance := float64(valueOfSteps) * stepLength
	distance = distance / mInKm

	calories, err := spentcalories.WalkingSpentCalories(valueOfSteps, weight, height, durationOfWalk)
	if err != nil {
		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", valueOfSteps, distance, calories)
	return result
}
