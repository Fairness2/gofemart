package orders

import (
	"encoding/json"
	database "gofemart/internal/databse"
	"gofemart/internal/helpers"
	"gofemart/internal/luna"
	"gofemart/internal/models"
	"gofemart/internal/ordercheck"
	"gofemart/internal/repositories"
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
		helpers.ProcessErrorWithStatus("Luna check failed", http.StatusUnprocessableEntity, response)
		return
	}

	// Берём авторизованного пользователя
	user, ok := request.Context().Value("user").(*models.User)
	if !ok {
		helpers.ProcessErrorWithStatus("User not found", http.StatusUnauthorized, response)
		return
	}

	rep := repositories.NewOrderRepository(request.Context(), database.DBx)

	order, ok, err := getOrderFromBd(rep, strBody)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}
	if ok {
		if order.UserId != user.Id {
			helpers.ProcessErrorWithStatus("order was loaded by another user", http.StatusConflict, response)
			return
		}
		if order.UserId == user.Id {
			helpers.ProcessErrorWithStatus("order was already loaded", http.StatusOK, response)
			return
		}
	}

	// Создаём новый ордер и отправляем его в проверочную
	order = models.NewOrder(strBody, user.Id)
	if err := saveOrder(rep, order); err != nil {
		helpers.SetInternalError(err, response)
		return
	}
	if _, err := sendToQueue(order); err != nil {
		helpers.SetInternalError(err, response)
		return
	}
	helpers.ProcessErrorWithStatus("new order number accepted for processing", http.StatusAccepted, response)
}

func getOrderFromBd(rep *repositories.OrderRepository, number string) (*models.Order, bool, error) {
	return rep.GetOrderByNumber(number)
}

func saveOrder(rep *repositories.OrderRepository, order *models.Order) error {
	return rep.CreateOrder(order)
}

func sendToQueue(order *models.Order) (bool, error) {
	return ordercheck.CheckPool.Push(order)
}

func GetOrdersHandler(response http.ResponseWriter, request *http.Request) {
	// Берём авторизованного пользователя
	user, ok := request.Context().Value("user").(*models.User)
	if !ok {
		helpers.ProcessErrorWithStatus("User not found", http.StatusUnauthorized, response)
		return
	}

	rep := repositories.NewOrderRepository(request.Context(), database.DBx)
	orders, err := rep.GetOrdersByUserWithAccrual(user.Id)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}

	if len(orders) == 0 {
		helpers.ProcessErrorWithStatus("no orders found", http.StatusNoContent, response)
		return
	}

	res, err := json.Marshal(orders)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}
	if err := helpers.SetHTTPResponse(response, http.StatusOK, res); err != nil {
		helpers.SetInternalError(err, response)
	}
}

func GetOrdersWwithdrawalsHandler(response http.ResponseWriter, request *http.Request) {
	// Берём авторизованного пользователя
	user, ok := request.Context().Value("user").(*models.User)
	if !ok {
		helpers.ProcessErrorWithStatus("User not found", http.StatusUnauthorized, response)
		return
	}

	rep := repositories.NewOrderRepository(request.Context(), database.DBx)
	orders, err := rep.GetOrdersByUserWithdraw(user.Id)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}

	if len(orders) == 0 {
		helpers.ProcessErrorWithStatus("no orders withdrawals", http.StatusNoContent, response)
		return
	}

	res, err := json.Marshal(orders)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}
	if err := helpers.SetHTTPResponse(response, http.StatusOK, res); err != nil {
		helpers.SetInternalError(err, response)
	}
}
