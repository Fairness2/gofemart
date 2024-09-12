package balance

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/asaskevich/govalidator"
	database "gofemart/internal/databse"
	"gofemart/internal/gofemarterrors"
	"gofemart/internal/helpers"
	"gofemart/internal/luna"
	"gofemart/internal/models"
	"gofemart/internal/payloads"
	"gofemart/internal/repositories"
	"gofemart/internal/services"
	"io"
	"net/http"
)

func WithdrawHandler(response http.ResponseWriter, request *http.Request) {
	// Читаем тело запроса
	body, err := getBody(request)
	if err != nil {
		helpers.ProcessRequestErrorWithBody(err, response)
		return
	}

	// Проверим полученный номер алгоритмом луна
	ok, err := luna.Check(body.OrderNumber)
	if err != nil {
		helpers.ProcessErrorWithStatus(err.Error(), http.StatusUnprocessableEntity, response)
		return
	}
	if !ok {
		helpers.ProcessErrorWithStatus("Luna check failed", http.StatusUnprocessableEntity, response)
		return
	}

	rep := repositories.NewOrderRepository(request.Context(), database.DBx)
	_, exists, err := rep.GetOrderByNumber(body.OrderNumber)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}
	if exists {
		helpers.ProcessErrorWithStatus("Order already exists", http.StatusUnprocessableEntity, response)
		return
	}

	accrualRep := repositories.NewAccountRepository(request.Context(), database.DBx)
	_, exists, err = accrualRep.GetWithdrawByOrder(body.OrderNumber)
	if err != nil {
		helpers.SetInternalError(err, response)
		return
	}
	if exists {
		helpers.ProcessErrorWithStatus("Withdraw already exists", http.StatusUnprocessableEntity, response)
		return
	}

	// Берём авторизованного пользователя
	user, ok := request.Context().Value("user").(*models.User)
	if !ok {
		helpers.ProcessErrorWithStatus("User not found", http.StatusUnauthorized, response)
		return
	}

	service := getBalanceService(request.Context())
	order := &models.Order{Number: body.OrderNumber}
	if err := service.Spend(user, body.Sum, order); err != nil {
		if errors.Is(err, services.ErrorNotEnoughItems) {
			helpers.ProcessErrorWithStatus(err.Error(), http.StatusPaymentRequired, response)
		} else {
			helpers.SetInternalError(err, response)
		}
		return
	}

	helpers.ProcessErrorWithStatus("Success", http.StatusOK, response)
}

// getBody получаем тело для регистрации
func getBody(request *http.Request) (*payloads.Withdraw, error) {
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

func getBalanceService(ctx context.Context) *services.BalanceService {
	return services.NewBalanceService(ctx)
}

func GetBalanceHandler(response http.ResponseWriter, request *http.Request) {
	// Берём авторизованного пользователя
	user, ok := request.Context().Value("user").(*models.User)
	if !ok {
		helpers.ProcessErrorWithStatus("User not found", http.StatusUnauthorized, response)
		return
	}
	rep := repositories.NewAccountRepository(request.Context(), database.DBx)
	bal, err := rep.GetBalance(user.Id)
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
