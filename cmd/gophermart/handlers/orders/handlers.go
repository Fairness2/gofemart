package orders

import (
	"encoding/json"
	"gofemart/internal/helpers"
	"gofemart/internal/luna"
	"gofemart/internal/models"
	"gofemart/internal/ordercheck"
	"gofemart/internal/repositories"
	"gofemart/internal/token"
	"io"
	"net/http"
	"strings"
)

type Handlers struct {
	dbPool repositories.SQLExecutor
}

func NewHandlers(dbPool repositories.SQLExecutor) *Handlers {
	return &Handlers{
		dbPool: dbPool,
	}
}

func (h *Handlers) RegisterOrderHandler(response http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}
	strBody := strings.Trim(string(body), " \n\r")

	// Проверим полученный номер алгоритмом луна
	ok, err := luna.Check(strBody)
	if err != nil {
		helpers.ProcessResponseWithStatus(err.Error(), http.StatusBadRequest, response)
		return
	}
	if !ok {
		helpers.ProcessResponseWithStatus("Luna check failed", http.StatusUnprocessableEntity, response)
		return
	}

	// Берём авторизованного пользователя
	user, ok := request.Context().Value(token.UserKey).(*models.User)
	if !ok {
		helpers.ProcessResponseWithStatus("User not found", http.StatusUnauthorized, response)
		return
	}

	rep := repositories.NewOrderRepository(request.Context(), h.dbPool)

	order, ok, err := h.getOrderFromBd(rep, strBody)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}
	if ok {
		if order.UserID != user.ID {
			helpers.ProcessResponseWithStatus("order was loaded by another user", http.StatusConflict, response)
			return
		}
		if order.UserID == user.ID {
			helpers.ProcessResponseWithStatus("order was already loaded", http.StatusOK, response)
			return
		}
	}

	// Создаём новый ордер и отправляем его в проверочную
	order = models.NewOrder(strBody, user.ID)
	if err := h.saveOrder(rep, order); err != nil {
		helpers.SetInternalError(err, response)
		return
	}
	if _, err := h.sendToQueue(order); err != nil {
		helpers.SetInternalError(err, response)
		return
	}
	helpers.ProcessResponseWithStatus("new order number accepted for processing", http.StatusAccepted, response)
}

func (h *Handlers) getOrderFromBd(rep *repositories.OrderRepository, number string) (*models.Order, bool, error) {
	return rep.GetOrderByNumber(number)
}

func (h *Handlers) saveOrder(rep *repositories.OrderRepository, order *models.Order) error {
	return rep.CreateOrder(order)
}

func (h *Handlers) sendToQueue(order *models.Order) (bool, error) {
	return ordercheck.CheckPool.Push(order)
}

func (h *Handlers) GetOrdersHandler(response http.ResponseWriter, request *http.Request) {
	// Берём авторизованного пользователя
	user, ok := request.Context().Value(token.UserKey).(*models.User)
	if !ok {
		helpers.ProcessResponseWithStatus("User not found", http.StatusUnauthorized, response)
		return
	}

	rep := repositories.NewOrderRepository(request.Context(), h.dbPool)
	orders, err := rep.GetOrdersByUserWithAccrual(user.ID)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}

	if len(orders) == 0 {
		helpers.ProcessResponseWithStatus("no orders found", http.StatusNoContent, response)
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

func (h *Handlers) GetOrdersWwithdrawalsHandler(response http.ResponseWriter, request *http.Request) {
	// Берём авторизованного пользователя
	user, ok := request.Context().Value(token.UserKey).(*models.User)
	if !ok {
		helpers.ProcessResponseWithStatus("User not found", http.StatusUnauthorized, response)
		return
	}

	rep := repositories.NewOrderRepository(request.Context(), h.dbPool)
	orders, err := rep.GetOrdersByUserWithdraw(user.ID)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}

	if len(orders) == 0 {
		helpers.ProcessResponseWithStatus("no orders withdrawals", http.StatusNoContent, response)
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
