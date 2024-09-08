package orders

import (
	"gofemart/internal/helpers"
	"gofemart/internal/luna"
	"gofemart/internal/models"
	"io"
	"net/http"
	"strings"
)

func RegisterOrderHandler(response http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}
	strBody := strings.Trim(string(body), " \n\r")

	// Проверим полученный номер алгоритмом луна
	ok, err := luna.Check(strBody)
	if err != nil {
		helpers.ProcessErrorWithStatus(err.Error(), http.StatusBadRequest, response)
		return
	}
	if !ok {
		helpers.ProcessErrorWithStatus("Luna check failed", http.StatusBadRequest, response)
		return
	}

	// Берём авторизованного пользователя
	user, ok := request.Context().Value("user").(*models.User)
	if !ok {
		helpers.ProcessErrorWithStatus("User not found", http.StatusUnauthorized, response)
		return
	}

	// Создаём новый ордер и отправляем его в проверочную
	_ = models.NewOrder(strBody, user.Id)

}
