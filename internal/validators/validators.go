package validators

import (
	"github.com/asaskevich/govalidator"
	"gofemart/internal/luna"
)

func init() {
	govalidator.CustomTypeTagMap.Set("luna", func(i interface{}, o interface{}) bool {
		num, ok := i.(string)
		// не удалось преобразовать ввод в строку
		if !ok {
			return false
		}
		res, err := luna.Check(num)
		return err == nil && res
	})
}
