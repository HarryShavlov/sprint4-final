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
	/*
		1. Разделить строку на слайс строк.
		2. Проверить, чтобы длина слайса была равна 3, так как в строке данных у нас количество шагов, вид активности и продолжительность.
		3. Преобразовать первый элемент слайса (количество шагов) в тип int. Обработать возможные ошибки. При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку.
		4. Преобразовать третий элемент слайса в time.Duration. В пакете time есть метод для парсинга строки в time.Duration. Обработать возможные ошибки. При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку.
		5. Если всё прошло без ошибок, верните количество шагов, вид активности, продолжительность и nil (для ошибки).
	*/
	dataSlice := strings.Split(data, ",")

	if len(dataSlice) != 3 {
		return 0, "", time.Duration(0), fmt.Errorf("некорректный формат")
	}

	steps, err := strconv.Atoi(dataSlice[0])
	if err != nil {
		return 0, "", time.Duration(0), err
	}

	if steps <= 0 {
		return 0, "", time.Duration(0), fmt.Errorf("неверныые шаги")
	}

	timeOfSteps, err := time.ParseDuration(dataSlice[2])
	if err != nil {
		return 0, "", time.Duration(0), err
	}

	if timeOfSteps <= 0 {
		return 0, "", time.Duration(0), fmt.Errorf("неверная продолжительность")
	}

	return steps, dataSlice[1], timeOfSteps, nil

}

func distance(steps int, height float64) float64 {
	/*
		1. рассчитайте длину шага. Для этого умножьте высоту пользователя на коэффициент длины шага stepLengthCoefficient. Соответствующая константа уже определена в пакете.
		2. умножьте пройденное количество шагов на длину шага.
		3. разделите полученное значение на число метров в километре (mInKm, константа определена в пакете).
	*/
	lenOfSteps := height * stepLengthCoefficient

	// Решил объединить перевод в км с подсчетом расстояния в метрах
	return lenOfSteps * float64(steps) / mInKm

}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	/*
		1. Проверить, что продолжительность duration больше 0. Если это не так, вернуть 0.
		2. Вычислить дистанцию с помощью distance().
		3. Вычислить и вернуть среднюю скорость. Для этого разделите дистанцию на продолжительность в часах. Чтобы перевести продолжительность в часы, воспользуйтесь функцией из пакета time.
	*/

	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)

	return dist / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {

	/*
		1. Получить значения из строки данных с помощью функции parseTraining(), обработать возможные ошибки и вывести их в лог с помощью log.Println(err).
		2. Проверить, какой вид тренировки был передан в строке, которую парсили (лучше использовать switch). Для каждого из видов тренировки рассчитать дистанцию, среднюю скорость и калории.
		3. Для каждого вида тренировки сформировать и вернуть строку, образец которой был представлен выше.
			Тип тренировки: Бег
			Длительность: 0.75 ч.
			Дистанция: 10.00 км.
			Скорость: 13.34 км/ч
			Сожгли калорий: 18621.75
		4. Если был передан неизвестный тип тренировки, вернуть ошибку с текстом неизвестный тип тренировки.
	*/
	steps, typeOfTrain, timeOfSteps, err := parseTraining(data)

	if err != nil {
		log.Println(err)
	}
	var calories float64 = 0

	dist := distance(steps, height)

	averageSpeed := meanSpeed(steps, height, timeOfSteps)

	switch typeOfTrain {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, timeOfSteps)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", typeOfTrain, timeOfSteps.Hours(), dist, averageSpeed, calories), nil

	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, timeOfSteps)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", typeOfTrain, timeOfSteps.Hours(), dist, averageSpeed, calories), nil
	}

	return "", fmt.Errorf("неизвестный тип тренировки")

}

func checkParametrs(steps int, weight, height float64, duration time.Duration) error {

	if duration <= 0 {
		return fmt.Errorf("неверная продолжительность")
	}

	if steps <= 0 {
		return fmt.Errorf("неверные шаги")
	}

	if weight <= 0 {
		return fmt.Errorf("неверный вес")
	}

	if height <= 0 {
		return fmt.Errorf("неверный рост")
	}

	return nil

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	/*
		1. Проверить входные параметры на корректность. Если параметры некорректны, вернуть 0 калорий и соответствующую ошибку.
		2. Рассчитать среднюю скорость с помощью meanSpeed().
		3. Рассчитать и вернуть количество калорий. Для этого:
			a. Переведите продолжительность в минуты с помощью функции из пакета time.
			b. Умножьте вес пользователя на среднюю скорость и продолжительность в минутах.
			c. Разделите результат на число минут в часе для получения количества потраченных калорий.
				(weight * meanSpeed * durationInMinutes) / minInH
	*/
	err := checkParametrs(steps, weight, height, duration)

	if err != nil {
		return 0, err
	}

	averageSpeed := meanSpeed(steps, height, duration)

	return (weight * averageSpeed * duration.Minutes()) / minInH, nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	/*
		1. Проверить входные параметры на корректность. Если параметры некорректны, вернуть 0 калорий и соответствующую ошибку.
		2. Рассчитать среднюю скорость с помощью meanSpeed().
		3. Рассчитать количество калорий. Для этого:
			a. Переведите продолжительность в минуты с помощью функции из пакета time.
			b. Умножьте вес пользователя на среднюю скорость и продолжительность в минутах.
			c. Разделите результат на число минут в часе для получения количества потраченных калорий.
			d. Умножить полученное число калорий на корректирующий коэффициент walkingCaloriesCoefficient. Соответствующая константа объявлена в пакете. Вернуть полученное значение.
	*/

	err := checkParametrs(steps, weight, height, duration)

	if err != nil {
		return 0, err
	}

	averageSpeed := meanSpeed(steps, height, duration)

	return (weight * averageSpeed * duration.Minutes()) / minInH * walkingCaloriesCoefficient, nil

}
