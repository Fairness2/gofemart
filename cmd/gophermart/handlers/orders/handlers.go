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

// Handlers Хэндлеры работы с заказами
type Handlers struct {
	dbPool repositories.SQLExecutor
}

// NewHandlers создает новый экземпляр Handlers с предоставленным SQLExecutor.
func NewHandlers(dbPool repositories.SQLExecutor) *Handlers {
	return &Handlers{
		dbPool: dbPool,
	}
}

// RegisterOrderHandler обрабатывает запрос на регистрацию заказа.
// Он считывает номер заказа, проверяет его с помощью алгоритма Луна
// @Summary Регистрирует новый заказ
// @Description обрабатывает запрос на регистрацию заказа.
// @Tags Заказы
// @Accept json
// @Produce json
// @Param order body string true "Order number"
// @Success 200 {object} payloads.ErrorResponseBody
// @Success 202 {object} payloads.ErrorResponseBody
// @Failure 400 {object} payloads.ErrorResponseBody
// @Failure 401 {object} payloads.ErrorResponseBody
// @Failure 409 {object} payloads.ErrorResponseBody
// @Failure 422 {object} payloads.ErrorResponseBody
// @Failure 500 {object} payloads.ErrorResponseBody
// @Router /api/user/orders [post]
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

// getOrderFromBd извлекает заказ из базы данных на основе предоставленного номера заказа.
// Он возвращает заказ, логическое значение, указывающее, был ли заказ найден, и любые ошибки, возникшие в ходе процесса.
func (h *Handlers) getOrderFromBd(rep *repositories.OrderRepository, number string) (*models.Order, bool, error) {
	return rep.GetOrderByNumber(number)
}

// saveOrder сохраняет новый заказ в базе данных, используя предоставленный репозиторий заказов.
func (h *Handlers) saveOrder(rep *repositories.OrderRepository, order *models.Order) error {
	return rep.CreateOrder(order)
}

// sendToQueue отправляет заказ в очередь для дальнейшей обработки.
// Возвращает логическое значение, указывающее на успешность операции, и ошибку, если она возникла.
func (h *Handlers) sendToQueue(order *models.Order) (bool, error) {
	return ordercheck.CheckPool.Push(order)
}

// GetOrdersHandler обрабатывает запросы на получение списка заказов для аутентифицированного пользователя.
// @Summary Получить список заказов
// @Description Возвращает список заказов для аутентифицированного пользователя.
// @Tags Заказы
// @Produce  json
// @Success 200 {array} models.OrderWithAccrual "Список заказов"
// @Failure 204 {string} payloads.ErrorResponseBody
// @Failure 401 {string} payloads.ErrorResponseBody
// @Failure 500 {string} payloads.ErrorResponseBody
// @Router /api/user/orders [get]
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

// GetOrdersWwithdrawalsHandler обрабатывает запросы на получение списка заказов со снятием средств для аутентифицированного пользователя.
// @Summary Получить заказы со снятием средств
// @Description Возвращает список заказов со снятием средств для аутентифицированного пользователя.
// @Tags Заказы
// @Produce  json
// @Success 200 {array} models.OrderWithdraw "Список заказов"
// @Failure 204 {string} payloads.ErrorResponseBody
// @Failure 401 {string} payloads.ErrorResponseBody
// @Failure 500 {string} payloads.ErrorResponseBody
// @Router /api/user/withdrawals [get]
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
