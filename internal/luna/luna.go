package luna

import (
	"errors"
	"strconv"
)

// ErrorIncorrectNumber ошибка, что числовая последовательность содержит недопустимые символы
var ErrorIncorrectNumber = errors.New("incorrect number")

// Check проверка строкового номера используя алгоритм Луна
// Если число содержит нечисловые символы, оно возвращает ErrorIncorrectNumber.
func Check(number string) (bool, error) {
	if number == "" {
		return false, ErrorIncorrectNumber
	}
	var sum int
	length := len(number)
	num := []rune(number)
	j := 0
	for i := length; i > 0; i-- {
		j++
		n, err := strconv.Atoi(string(num[i-1]))
		if err != nil {
			return false, ErrorIncorrectNumber
		}
		if j%2 != 0 {
			sum += n
			continue
		}
		n = n * 2
		if n > 9 {
			n -= 9
		}
		sum += n
	}

	return sum%10 == 0, nil
}
