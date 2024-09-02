package validators

import (
	"github.com/asaskevich/govalidator"
	"gofemart/internal/logger"
	"gofemart/internal/repositories"
)

func init() {
	govalidator.CustomTypeTagMap.Set("uniqueLogin", func(i interface{}, o interface{}) bool {
		res, err := repositories.UserR.UserExists(i.(string))
		if err != nil { // Если нашли ошибку, то записываем еёб жаль нельзя закрыть запрос(
			logger.Log.Error(err)
			return false
		}
		return res
	})
}
