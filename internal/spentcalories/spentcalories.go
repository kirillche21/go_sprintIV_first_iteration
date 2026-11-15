package spentcalories

import (
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

func parseTraining(data string) (int, string, time.Duration, error) {
	split := strings.Split(data, ",")

	if len(split) != 3 {
		return 0, "", 0, fmt.Errorf("Ошибка: Длина слайса %d != 3", len(split))
	}

	valueOfSteps, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Ошибка преобразования 1-го элемента: %s", err)
	}
	if valueOfSteps <= 0 {
		return 0, "", 0, fmt.Errorf("Ошибка: количество шагов < 0 %s", err)
	}

	typeOfActivity := split[1]

	durationOfWalk, err := time.ParseDuration(split[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Ошибка преобразования 3-го элемента: %s", err)
	}
	if durationOfWalk <= 0 {
		return 0, "", 0, fmt.Errorf("Ошибка: длительность прогулки < 0 %d", durationOfWalk)
	}

	return valueOfSteps, typeOfActivity, durationOfWalk, nil
}

func distance(steps int, height float64) float64 {
	lenOfStep := height * stepLengthCoefficient

	distanceCovered := lenOfStep * float64(steps)

	distanceInKm := distanceCovered / float64(mInKm)

	return distanceInKm

}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distance := distance(steps, height)

	averageSpeed := distance / (float64(duration) / float64(time.Hour))

	return averageSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	textOfError := "неизвестный тип тренировки"

	valueOfSteps, typeOfActivity, durationOfWalk, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	switch typeOfActivity {
	case "Бег":
		d := distance(valueOfSteps, height)
		s := meanSpeed(valueOfSteps, height, durationOfWalk)
		c, err := RunningSpentCalories(valueOfSteps, weight, height, durationOfWalk)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", typeOfActivity, durationOfWalk.Hours(), d, s, c), nil
	case "Ходьба":
		d := distance(valueOfSteps, height)
		s := meanSpeed(valueOfSteps, height, durationOfWalk)
		c, err := WalkingSpentCalories(valueOfSteps, weight, height, durationOfWalk)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", typeOfActivity, durationOfWalk.Hours(), d, s, c), nil
	default:
		return "", fmt.Errorf("%s", textOfError)
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	textOfError := "Ошибка: Вводные данные некорректны"
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("%s", textOfError)
	}

	averageSpeed := meanSpeed(steps, height, duration)

	durationInMin := duration.Minutes()

	caloriesBurned := (weight * averageSpeed * durationInMin) / float64(minInH)

	return caloriesBurned, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	textOfError := "Ошибка: Вводные данные некорректны"
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("%s", textOfError)
	}

	averageSpeed := meanSpeed(steps, height, duration)

	durationInMin := float64(duration / time.Minute)

	caloriesBurned := ((weight * averageSpeed * durationInMin) / float64(minInH)) * walkingCaloriesCoefficient

	return caloriesBurned, nil
}
