package balance

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/asaskevich/govalidator"
	"gofemart/internal/gofemarterrors"
	"gofemart/internal/helpers"
	"gofemart/internal/luna"
	"gofemart/internal/models"
	"gofemart/internal/payloads"
	"gofemart/internal/repositories"
	"gofemart/internal/services"
	"gofemart/internal/token"
	"io"
	"net/http"
)

// Handlers для обработки запросов, связанных с балансом.
type Handlers struct {
	dbPool repositories.SQLExecutor
}

// NewHandlers инициализирует и возвращает новый экземпляр Handlers с предоставленным dbPool. SQLExecutor.
func NewHandlers(dbPool repositories.SQLExecutor) *Handlers {
	return &Handlers{
		dbPool: dbPool,
	}
}

// WithdrawHandler обрабатывает HTTP-запрос на вывод указанной суммы с баланса счета пользователя.
// Он проверяет номер заказа с помощью алгоритма Luna, проверяет, существует ли уже заказ или вывод,
// извлекает аутентифицированного пользователя и пытается потратить указанную сумму.
func (b *Handlers) WithdrawHandler(response http.ResponseWriter, request *http.Request) {
	// Читаем тело запроса
	body, err := b.getBody(request)
	if err != nil {
		helpers.ProcessRequestErrorWithBody(err, response)
		return
	}

	// Проверим полученный номер алгоритмом луна
	ok, err := luna.Check(body.OrderNumber)
	if err != nil {
		helpers.ProcessResponseWithStatus(err.Error(), http.StatusUnprocessableEntity, response)
		return
	}
	if !ok {
		helpers.ProcessResponseWithStatus("Luna check failed", http.StatusUnprocessableEntity, response)
		return
	}

	rep := repositories.NewOrderRepository(request.Context(), b.dbPool)
	_, exists, err := rep.GetOrderByNumber(body.OrderNumber)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}
	if exists {
		helpers.ProcessResponseWithStatus("Order already exists", http.StatusUnprocessableEntity, response)
		return
	}

	accrualRep := repositories.NewAccountRepository(request.Context(), b.dbPool)
	_, exists, err = accrualRep.GetWithdrawByOrder(body.OrderNumber)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}
	if exists {
		helpers.ProcessResponseWithStatus("Withdraw already exists", http.StatusUnprocessableEntity, response)
		return
	}

	// Берём авторизованного пользователя
	user, ok := request.Context().Value(token.UserKey).(*models.User)
	if !ok {
		helpers.ProcessResponseWithStatus("User not found", http.StatusUnauthorized, response)
		return
	}

	service := b.getBalanceService(request.Context())
	order := &models.Order{Number: body.OrderNumber}
	if err := service.Spend(user, body.Sum, order); err != nil {
		if errors.Is(err, services.ErrorNotEnoughItems) {
			helpers.ProcessResponseWithStatus(err.Error(), http.StatusPaymentRequired, response)
		} else {
			helpers.SetInternalError(err, response)
		}
		return
	}

	helpers.ProcessResponseWithStatus("Success", http.StatusOK, response)
}

// getBody получаем тело для регистрации
func (b *Handlers) getBody(request *http.Request) (*payloads.Withdraw, error) {
	// Читаем тело запроса
	rawBody, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	// Парсим тело в структуру запроса
	var body payloads.Withdraw
	err = json.Unmarshal(rawBody, &body)
	if err != nil {
		return nil, &gofemarterrors.RequestError{InternalError: err, HTTPStatus: http.StatusBadRequest}
	}

	result, err := govalidator.ValidateStruct(body)
	if err != nil {
		return nil, err
	}

	if !result {
		return nil, &gofemarterrors.RequestError{InternalError: errors.New("bad request"), HTTPStatus: http.StatusBadRequest}
	}

	return &body, nil
}

// getBalanceService создает и возвращает новый экземпляр BalanceService, используя предоставленный контекст и пул базы данных.
func (b *Handlers) getBalanceService(ctx context.Context) *services.BalanceService {
	return services.NewBalanceService(ctx, b.dbPool)
}

// GetBalanceHandler обрабатывает HTTP-запросы для получения баланса счета аутентифицированного пользователя.
func (b *Handlers) GetBalanceHandler(response http.ResponseWriter, request *http.Request) {
	// Берём авторизованного пользователя
	user, ok := request.Context().Value(token.UserKey).(*models.User)
	if !ok {
		helpers.ProcessResponseWithStatus("User not found", http.StatusUnauthorized, response)
		return
	}
	rep := repositories.NewAccountRepository(request.Context(), b.dbPool)
	bal, err := rep.GetBalance(user.ID)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}

	res, err := json.Marshal(bal)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}
	if err := helpers.SetHTTPResponse(response, http.StatusOK, res); err != nil {
		helpers.SetInternalError(err, response)
	}
}
