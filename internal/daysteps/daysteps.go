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
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	/*
		1. Разделить строку на слайс строк.
		2. Проверить, чтобы длина слайса была равна 2, так как в строке данных у нас количество шагов и продолжительность.
		3. Преобразовать первый элемент слайса (количество шагов) в тип int. Обработать возможные ошибки. При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку.
		4. Проверить: количество шагов должно быть больше 0. Если это не так, вернуть нули и ошибку.
		5. Преобразовать второй элемент слайса в time.Duration. В пакете time есть метод для парсинга строки в time.Duration. Обработать возможные ошибки. При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку.
		6. Если всё прошло без ошибок, верните количество шагов, продолжительность и nil (для ошибки).
	*/
	dataSlice := strings.Split(data, ",")
	if len(dataSlice) != 2 {
		return 0, time.Duration(0), fmt.Errorf("2 parameters are expected in data")
	}

	steps, err := strconv.Atoi(dataSlice[0])
	if err != nil {
		return 0, time.Duration(0), err
	}

	if steps <= 0 {
		return 0, time.Duration(0), fmt.Errorf("incorrect number of steps")
	}

	timeOfSteps, err := time.ParseDuration(dataSlice[1])
	if err != nil {
		return 0, time.Duration(0), err
	}

	return steps, timeOfSteps, nil

}

func DayActionInfo(data string, weight, height float64) string {
	/*
		1. Получить данные о количестве шагов и продолжительности прогулки с помощью функции parsePackage(). В случае возникновения ошибки вывести её на экран и вернуть пустую строку.
		2. Проверить, чтобы количество шагов было больше 0. В противном случае вернуть пустую строку.
		3. Вычислить дистанцию в метрах. Дистанция равна произведению количества шагов на длину шага. Константа stepLength (длина шага) уже определена в коде.
		4. Перевести дистанцию в километры, разделив её на число метров в километре (константа mInKm, определена в пакете).
		5. Вычислить количество калорий, потраченных на прогулке. Функция для вычисления калорий WalkingSpentCalories() будет определена в пакете spentcalories, которую вы тоже реализуете.
		6. Сформировать строку, которую будете возвращать, пример которой был представлен выше.(Количество шагов: 792. \n Дистанция составила 0.51 км. \n Вы сожгли 221.33 ккал. )
	*/
	steps, timeOfSteps, err := parsePackage(data)

	// Подумал, что если возвращаемое значение одинаковое, то могу в 1 if уместить
	if err != nil || steps == 0 {
		return ""
	}
	// if steps == 0 {
	// 	return ""
	// }

	// Дистанцию сразу в км посчитал
	dist := float64(steps) * stepLength / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, timeOfSteps)

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, dist, calories)
}
