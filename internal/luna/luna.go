package luna

import (
	"errors"
	"strconv"
)

var ErrorIncorrectNumber = errors.New("incorrect number")

func Check(number string) (bool, error) {
	var sum int
	/*if length := len(number); length%2 == 1 {
		return false, ErrorIncorrectNumber
	}*/
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
